package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/interfaces/repository"
	"ml-prediction/internal/app/interfaces/usecase"
	"ml-prediction/internal/app/model"
	"os/exec"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type customerUsecase struct {
	custPredRepo repository.CustomerRepository
}

func NewcustomerUsecase(custPredRepo repository.CustomerRepository) usecase.CustomerUsecase {
	return &customerUsecase{custPredRepo}
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

	// Run Python script
	cmd := exec.Command("python3", "scripts/model.py")
	cmd.Stdin = bytes.NewReader(inputJSON)

	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, errors.New("User dengan CIF yang diberikan telah ada!")
	}

	// Parse Python output
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
	var TopProduk []string
	for k, v := range predictions {
		TopProduk = append(TopProduk, k)
		sortedPreds = append(sortedPreds, kv{k, v})
	}

	sort.Slice(sortedPreds, func(i, j int) bool {
		return sortedPreds[i].Value > sortedPreds[j].Value
	})

	combined := &model.Customer{
		Nama:               req.Nama,
		NamaPerusahaan:     req.NamaPerusahaan,
		Umur:               req.Umur,
		Penghasilan:        req.Penghasilan,
		Payroll:            req.Payroll,
		Gender:             req.Gender,
		StatusPerkawinan:   req.StatusPerkawinan,
		Segmen:             req.Segmen,
		ProdukEksisting:    req.ProdukEksisting,
		AktivitasTransaksi: req.AktivitasTransaksi,
		TopProduk:          pq.StringArray(TopProduk),
	}

	data, err := s.custPredRepo.Create(c, combined)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Gagal menambahkan data prediksi: %v", err))
	}

	return data, nil
}
