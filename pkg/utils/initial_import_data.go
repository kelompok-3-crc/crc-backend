package utils

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"ml-prediction/pkg/helper"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

const maxConcurrentWorkers = 30

type workerResult struct {
	customer    model.Customer
	sortedPreds []predictionPair
	lineNum     int
	err         error
}

type predictionPair struct {
	Key   string
	Value float64
}

func ImportInitialCustomerData(ctx context.Context, db *gorm.DB) error {
	startTime := time.Now()

	var count int64
	if err := db.Model(&model.Customer{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check customer count: %v", err)
	}

	if count > 0 {
		log.Println("Customer data already exists, skipping import")
		return nil
	}

	projectRoot, err := helper.GetProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to get project root: %v", err)
	}

	dataPath := filepath.Join(projectRoot, "data.csv")

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {

		altDataPath := filepath.Join(projectRoot, "data", "data.csv")
		if _, err := os.Stat(altDataPath); os.IsNotExist(err) {
			return fmt.Errorf("CSV data file not found at %s or %s", dataPath, altDataPath)
		}
		dataPath = altDataPath
	}

	file, err := os.Open(dataPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	productMap := make(map[string]model.Product)
	var products []model.Product
	if err := tx.Find(&products).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch products: %v", err)
	}

	for _, p := range products {
		productMap[strings.ToLower(p.Prediksi)] = p
	}

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to read header: %v", err)
	}

	records := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error reading CSV: %v", err)
		}
		records = append(records, record)
	}

	totalRecords := len(records)
	log.Printf("Starting parallel processing of %d records with %d workers", totalRecords, maxConcurrentWorkers)

	jobs := make(chan struct {
		record  []string
		lineNum int
	}, totalRecords)
	results := make(chan workerResult, totalRecords)

	var wg sync.WaitGroup
	workerCount := maxConcurrentWorkers
	if workerCount > totalRecords {
		workerCount = totalRecords
	}

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg, projectRoot)
	}

	for i, record := range records {
		jobs <- struct {
			record  []string
			lineNum int
		}{record, i + 1}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	successCount := 0
	errorCount := 0

	processedResults := make([]workerResult, 0, totalRecords)
	for result := range results {
		if result.err != nil {
			log.Printf("Error processing record %d: %v", result.lineNum, result.err)
			errorCount++
			continue
		}
		processedResults = append(processedResults, result)
		successCount++

		if successCount%20 == 0 {
			log.Printf("Processed %d/%d records... (%.2f%%)",
				successCount, totalRecords, float64(successCount)/float64(totalRecords)*100)
		}
	}

	for _, result := range processedResults {

		customerWithoutProducts := result.customer
		customerWithoutProducts.CustomerProduk = nil
		if err := tx.Create(&customerWithoutProducts).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error creating customer at line %d: %v", result.lineNum, err)
		}

		for i := 0; i < len(result.sortedPreds); i++ {
			prodName := result.sortedPreds[i].Key
			if result.sortedPreds[i].Value == 0 {
				continue
			}

			product, exists := productMap[strings.ToLower(prodName)]
			if !exists {
				log.Printf("Warning: Product '%s' not found in database", prodName)
				continue
			}
			plafond := helper.CalculatePlafond(product.Prediksi, int64(customerWithoutProducts.Umur), customerWithoutProducts.Penghasilan, customerWithoutProducts.Payroll)

			customerProd := &model.CustomerProduct{
				CustomerID: customerWithoutProducts.Id,
				ProductID:  product.ID,
				Order:      i + 1,
				PlafonMin:  &plafond.MinPlafon,
				PlafonMax:  &plafond.MaxPlafon,
				TenorMin:   &plafond.MinTenor,
				TenorMax:   &plafond.MaxTenor,
			}

			if err := tx.Create(customerProd).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("error creating customer product at line %d: %v", result.lineNum, err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("Successfully imported %d customers from CSV with predictions (skipped %d) in %v",
		successCount, errorCount, duration)
	return nil
}

func worker(id int, jobs <-chan struct {
	record  []string
	lineNum int
}, results chan<- workerResult, wg *sync.WaitGroup, projectRoot string) {
	defer wg.Done()

	for job := range jobs {
		record := job.record
		lineNum := job.lineNum

		customer, err := parseCSVRow(record, lineNum)
		if err != nil {
			results <- workerResult{lineNum: lineNum, err: fmt.Errorf("error parsing CSV row: %v", err)}
			continue
		}

		mlInputData := &dto.PredictionRequest{
			CIF:                customer.CIF,
			Nama:               customer.Nama,
			NamaPerusahaan:     customer.NamaPerusahaan,
			NomorRekening:      customer.NomorRekening,
			NomorHp:            customer.NomorHp,
			Umur:               customer.Umur,
			Penghasilan:        customer.Penghasilan,
			Payroll:            customer.Payroll,
			Gender:             customer.Gender,
			StatusPerkawinan:   customer.StatusPerkawinan,
			Segmen:             customer.Segmen,
			ProdukEksisting:    customer.ProdukEksisting,
			AktivitasTransaksi: customer.AktivitasTransaksi,
		}
		if mlInputData.ProdukEksisting == nil {
			mlInputData.ProdukEksisting = pq.StringArray{}
		}

		inputJSON, err := json.Marshal(mlInputData)
		if err != nil {
			results <- workerResult{lineNum: lineNum, err: fmt.Errorf("error marshaling data: %v", err)}
			continue
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
			results <- workerResult{
				lineNum: lineNum,
				err:     fmt.Errorf("failed to run Python script: %v, %s", err, stderr.String()),
			}
			continue
		}

		var predictions map[string]float64
		err = json.Unmarshal(out.Bytes(), &predictions)
		if err != nil {
			results <- workerResult{lineNum: lineNum, err: fmt.Errorf("failed to parse predictions: %v", err)}
			continue
		}

		var sortedPreds []predictionPair
		for k, v := range predictions {
			sortedPreds = append(sortedPreds, predictionPair{k, v})
		}

		sort.Slice(sortedPreds, func(i, j int) bool {
			return sortedPreds[i].Value > sortedPreds[j].Value
		})

		results <- workerResult{
			customer:    customer,
			sortedPreds: sortedPreds,
			lineNum:     lineNum,
			err:         nil,
		}
	}
}

func parseCSVRow(record []string, lineNum int) (model.Customer, error) {
	if len(record) < 14 {
		return model.Customer{}, fmt.Errorf("invalid record length, expected 14 got %d", len(record))
	}

	cif := fmt.Sprintf("CIF%08d", lineNum)

	rekening := fmt.Sprintf("REK%010d", lineNum)

	namaPerusahaan := "PT Sample " + strconv.Itoa(lineNum)

	var produkEksisting pq.StringArray
	if record[13] != "" && record[13] != "nan" {
		prodExisting := strings.Split(record[13], ",")
		for i, prod := range prodExisting {
			prodExisting[i] = strings.TrimSpace(prod)
		}
		produkEksisting = prodExisting
	}

	gender := "MALE"
	if record[4] == "FEMALE" {
		gender = "FEMALE"
	}

	umur, err := strconv.Atoi(record[5])
	if err != nil {
		return model.Customer{}, fmt.Errorf("invalid age value: %v", err)
	}

	penghasilan, err := strconv.ParseInt(record[6], 10, 64)
	if err != nil {
		return model.Customer{}, fmt.Errorf("invalid income value: %v", err)
	}

	statusPerkawinan := false
	if record[7] == "Married" {
		statusPerkawinan = true
	}

	payroll := false
	if record[12] == "1" {
		payroll = true
	}

	customer := model.Customer{
		CIF:                cif,
		Nama:               "Customer " + strconv.Itoa(lineNum),
		NomorRekening:      rekening,
		NamaPerusahaan:     namaPerusahaan,
		ProdukEksisting:    produkEksisting,
		AktivitasTransaksi: record[8],
		NomorHp:            fmt.Sprintf("08%09d", lineNum),
		Segmen:             record[3],
		Address:            "Jl. Sample " + strconv.Itoa(lineNum),
		Job:                "Profession " + strconv.Itoa(lineNum%10+1),
		Penghasilan:        penghasilan,
		Umur:               umur,
		Gender:             gender,
		StatusPerkawinan:   statusPerkawinan,
		Payroll:            payroll,
	}

	return customer, nil
}
