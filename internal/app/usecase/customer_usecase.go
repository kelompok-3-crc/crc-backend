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
	// Validate unique fields
	if err := s.validateUniqueCustomerFields(c.Context(), req); err != nil {
		return nil, err
	}

	// Continue with existing code for prediction and customer creation
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

	// Use a different loop approach to ensure we get 3 valid products
	for i := 0; i < len(sortedPreds); i++ {
		// Skip products with prediction value of exactly 0
		fmt.Println("Prediksi:", sortedPreds[i].Key, "Nilai:", sortedPreds[i].Value)
		if sortedPreds[i].Value == 0 {
			continue
		}

		produk, err := s.produkRepo.FindByPrediksi(sortedPreds[i].Key)
		if err != nil {
			tx.Rollback()
			return nil, errors.New(fmt.Sprintf("Gagal menemukan produk: %v", err))
		}

		plafond := helper.CalculatePlafond(produk.Prediksi, int64(data.Umur), data.Penghasilan, data.Payroll)

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

	var fullCustomer model.Customer
	if err := tx.Preload("CustomerProduk").Where("id = ?", data.Id).First(&fullCustomer).Error; err != nil {
		tx.Rollback()
		return nil, errors.New(fmt.Sprintf("Gagal mengambil data lengkap customer: %v", err))
	}

	tx.Commit()
	return &fullCustomer, nil
}

// Add this new method for validation
func (s *customerUsecase) validateUniqueCustomerFields(ctx context.Context, req dto.PredictionRequest) error {
	// Check CIF uniqueness
	var cifCount int64
	if err := s.db.WithContext(ctx).Model(&model.Customer{}).
		Where("cif = ? AND deleted_at IS NULL", req.CIF).
		Count(&cifCount).Error; err != nil {
		return fmt.Errorf("error checking CIF uniqueness: %v", err)
	}
	if cifCount > 0 {
		return errors.New("customer dengan CIF tersebut sudah terdaftar")
	}

	// Check NomorRekening uniqueness (if provided)
	if req.NomorRekening != "" {
		var rekCount int64
		if err := s.db.WithContext(ctx).Model(&model.Customer{}).
			Where("nomor_rekening = ? AND deleted_at IS NULL", req.NomorRekening).
			Count(&rekCount).Error; err != nil {
			return fmt.Errorf("error checking Nomor Rekening uniqueness: %v", err)
		}
		if rekCount > 0 {
			return errors.New("customer dengan Nomor Rekening tersebut sudah terdaftar")
		}
	}

	// Check email uniqueness (if provided and if it's intended to be unique)
	if req.Email != "" {
		var emailCount int64
		if err := s.db.WithContext(ctx).Model(&model.Customer{}).
			Where("email = ? AND deleted_at IS NULL", req.Email).
			Count(&emailCount).Error; err != nil {
			return fmt.Errorf("error checking Email uniqueness: %v", err)
		}
		if emailCount > 0 {
			return errors.New("customer dengan Email tersebut sudah terdaftar")
		}
	}

	// Check phone number uniqueness (if provided and if it's intended to be unique)
	if req.NomorHp != "" {
		var phoneCount int64
		if err := s.db.WithContext(ctx).Model(&model.Customer{}).
			Where("nomor_hp = ? AND deleted_at IS NULL", req.NomorHp).
			Count(&phoneCount).Error; err != nil {
			return fmt.Errorf("error checking Nomor HP uniqueness: %v", err)
		}
		if phoneCount > 0 {
			return errors.New("customer dengan Nomor HP tersebut sudah terdaftar")
		}
	}

	return nil
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
