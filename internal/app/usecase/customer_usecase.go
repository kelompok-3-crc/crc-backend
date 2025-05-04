package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"
	"ml-prediction/internal/app/repository"
	"ml-prediction/pkg/helper"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CustomerUsecase interface {
	Create(c *fiber.Ctx, req dto.PredictionRequest) (*model.Customer, error)
	GetNewCustomers(ctx context.Context, NIP string, req *dto.CustomerSearchRequest) ([]dto.Customer, *dto.Pagination, error)
	GetAssignedCustomers(ctx context.Context, NIP string, req *dto.AssignedCustomerRequest) ([]dto.Customer, *dto.Pagination, error)
	GetCustomerDetail(ctx context.Context, NIP string, customerID string) (*dto.Customer, error)
}
type customerUsecase struct {
	custPredRepo repository.CustomerRepository
	userRepo     repository.UserRepository
	produkRepo   repository.ProductRepository
	db           *gorm.DB
}

func NewcustomerUsecase(custPredRepo repository.CustomerRepository, userRepo repository.UserRepository, produkRepo repository.ProductRepository, db *gorm.DB) CustomerUsecase {
	return &customerUsecase{custPredRepo, userRepo, produkRepo, db}
}
func (s *customerUsecase) Create(c *fiber.Ctx, req dto.PredictionRequest) (*model.Customer, error) {

	exists, err := s.custPredRepo.ExistsByCif(c, req.CIF)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("User dengan CIF yang diberikan telah ada!")
	}

	inputJSON, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("Gagal memproses data input!")
	}

	projectRoot, err := helper.GetProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to get project root: %v", err)
	}
	pythonPath := filepath.Join(projectRoot, "venv", "bin", "python3.11")
	scriptPath := filepath.Join(projectRoot, "scripts", "modelling.py")

	cmd := exec.Command(pythonPath, scriptPath)
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PYTHONPATH=%s", filepath.Join(projectRoot, "venv/lib/python3.11/site-packages")))
	cmd.Stdin = bytes.NewReader(inputJSON)

	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gagal menjalankan script Python: %v, %s", err, stderr.String()))
	}

	var predictions map[string]float64
	err = json.Unmarshal(out.Bytes(), &predictions)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse predictions: %v", err))
	}
	type kv struct {
		Key   string
		Value float64
	}

	var sortedPreds []kv
	for k, v := range predictions {
		sortedPreds = append(sortedPreds, kv{k, v})
	}

	sort.Slice(sortedPreds, func(i, j int) bool {
		return sortedPreds[i].Value > sortedPreds[j].Value
	})

	combined := &model.Customer{
		CIF:                req.CIF,
		Nama:               req.Nama,
		NamaPerusahaan:     req.NamaPerusahaan,
		NomorRekening:      req.NomorRekening,
		NomorHp:            req.NomorHp,
		Umur:               req.Umur,
		Penghasilan:        req.Penghasilan,
		Payroll:            req.Payroll,
		Gender:             req.Gender,
		StatusPerkawinan:   req.StatusPerkawinan,
		Segmen:             req.Segmen,
		ProdukEksisting:    req.ProdukEksisting,
		AktivitasTransaksi: req.AktivitasTransaksi,
		Email:              req.Email,
		Address:            req.Alamat,
		Job:                req.Pekerjaan,
	}

	tx := s.db.WithContext(c.Context()).Begin()
	customerWithoutProducts := *combined
	customerWithoutProducts.CustomerProduk = nil
	data, err := s.custPredRepo.CreateTx(tx, &customerWithoutProducts)
	if err != nil {
		tx.Rollback()
		return nil, errors.New(fmt.Sprintf("Gagal menambahkan data customer: %v", err))
	}
	customerProduct := []*model.CustomerProduct{}
	for i := 0; i < 3 && i < len(sortedPreds); i++ {
		produk, err := s.produkRepo.FindByPrediksi(sortedPreds[i].Key)
		if err != nil {
			tx.Rollback()
			return nil, errors.New(fmt.Sprintf("Gagal menemukan produk: %v", err))
		}

		plafond := CalculatePlafond(produk.Prediksi, int64(data.Umur), data.Penghasilan, data.Payroll)
		// Create a product relationship with customer ID
		customerProd := &model.CustomerProduct{
			CustomerID: data.Id,
			ProductID:  produk.ID,
			Order:      int(i + 1),
			PlafonMin:  &plafond.MinPlafon,
			PlafonMax:  &plafond.MaxPlafon,
			TenorMin:   &plafond.MinTenor,
			TenorMax:   &plafond.MaxTenor,
		}

		if err := tx.Create(customerProd).Error; err != nil {
			tx.Rollback()
			return nil, errors.New(fmt.Sprintf("Gagal menyimpan produk nasabah: %v", err))
		}

		customerProduct = append(customerProduct, customerProd)
	}

	// Fetch the complete customer with products
	var fullCustomer model.Customer
	if err := tx.Preload("CustomerProduk").Where("id = ?", data.Id).First(&fullCustomer).Error; err != nil {
		tx.Rollback()
		return nil, errors.New(fmt.Sprintf("Gagal mengambil data lengkap customer: %v", err))
	}

	tx.Commit()
	return &fullCustomer, nil
}

func (u *customerUsecase) GetNewCustomers(ctx context.Context, NIP string, req *dto.CustomerSearchRequest) ([]dto.Customer, *dto.Pagination, error) {

	user, err := u.userRepo.FindByNIP(NIP)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting user data: %v", err)
	}

	if user.Role != "marketing" {
		return nil, nil, errors.New("unauthorized access")
	}

	return u.custPredRepo.GetNewCustomers(req)
}

func (u *customerUsecase) GetAssignedCustomers(ctx context.Context, NIP string, req *dto.AssignedCustomerRequest) ([]dto.Customer, *dto.Pagination, error) {

	user, err := u.userRepo.FindByNIP(NIP)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting user data: %v", err)
	}

	if user.Role != "marketing" {
		return nil, nil, errors.New("unauthorized access")
	}

	return u.custPredRepo.GetAssignedCustomers(user.ID, req)
}

func (u *customerUsecase) GetCustomerDetail(ctx context.Context, NIP string, customerID string) (*dto.Customer, error) {

	user, err := u.userRepo.FindByNIP(NIP)
	if err != nil {
		return nil, fmt.Errorf("error getting user data: %v", err)
	}

	if user.Role != "marketing" {
		return nil, errors.New("unauthorized access")
	}

	return u.custPredRepo.GetCustomerDetail(user.ID, customerID)
}

type Plafond struct {
	MinPlafon uint64
	MaxPlafon uint64
	MinTenor  int
	MaxTenor  int
}

func CalculatePlafond(produk string, umur, penghasilan int64, payroll bool) Plafond {
	if produk == "mitraguna" {
		penghasilan := float64(penghasilan)
		minPlafon := 0.7 * penghasilan * 12
		maxPlafon := 0.7 * penghasilan * 15 * 12
		minTenor := 12
		maxTenor := 15 * 12

		return Plafond{
			MinPlafon: uint64(minPlafon),
			MaxPlafon: uint64(maxPlafon),
			MinTenor:  minTenor,
			MaxTenor:  maxTenor,
		}
	}

	if produk == "pensiun" {
		const PriceAkhirPensiun = 1000
		tenorMaks := (75 - int(umur)) * 12
		angsuranMaks := 0.9 * float64(penghasilan)
		plafondMaks := angsuranMaks * PriceAkhirPensiun

		return Plafond{
			MinPlafon: 0,
			MaxPlafon: uint64(plafondMaks),
			MinTenor:  0,
			MaxTenor:  tenorMaks,
		}
	}

	if produk == "griya" {
		const DSR = 0.4
		const MARGIN = 0.1
		const MIN_TENOR = 12
		const MAX_TENOR = 300

		angsuranMaks := float64(penghasilan) * DSR

		tenorBulanMax := float64(MAX_TENOR)
		plafonMax := (angsuranMaks * tenorBulanMax) / (1 + (MARGIN * tenorBulanMax / 12))

		tenorBulanMin := float64(MIN_TENOR)
		plafonMin := (angsuranMaks * tenorBulanMin) / (1 + (MARGIN * tenorBulanMin / 12))

		return Plafond{
			MinPlafon: uint64(plafonMin),
			MaxPlafon: uint64(plafonMax),
			MinTenor:  MIN_TENOR,
			MaxTenor:  MAX_TENOR,
		}
	}

	if produk == "oto" {
		const DSR = 0.5
		const MARGIN = 0.075
		const MIN_TENOR = 12
		const MAX_TENOR = 60

		angsuranMaks := float64(penghasilan) * DSR

		tenorBulanMax := float64(MAX_TENOR)
		pokokPembiayaanMax := (angsuranMaks * tenorBulanMax) / (1 + (MARGIN * tenorBulanMax / 12))

		tenorBulanMin := float64(MIN_TENOR)
		pokokPembiayaanMin := (angsuranMaks * tenorBulanMin) / (1 + (MARGIN * tenorBulanMin / 12))

		return Plafond{
			MinPlafon: uint64(pokokPembiayaanMin),
			MaxPlafon: uint64(pokokPembiayaanMax),
			MinTenor:  MIN_TENOR,
			MaxTenor:  MAX_TENOR,
		}
	}
	if produk == "hasanahcard" {
		var faktorLimit float64

		if umur < 35 {
			faktorLimit = 2.5
		} else if umur >= 35 && umur < 45 {

			faktorLimit = 2.0
		} else if umur >= 45 && umur < 60 {

			faktorLimit = 1.5
		} else {
			faktorLimit = 1.0
		}

		plafonMax := float64(penghasilan) * faktorLimit

		plafonMin := float64(penghasilan) * 1.0

		return Plafond{
			MinPlafon: uint64(plafonMin),
			MaxPlafon: uint64(plafonMax),
			MinTenor:  0,
			MaxTenor:  0,
		}
	}
	return Plafond{
		MinPlafon: 0,
		MaxPlafon: 0,
		MinTenor:  0,
		MaxTenor:  0,
	}
}
