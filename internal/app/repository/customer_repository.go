package repository

import (
	"errors"
	"fmt"
	"math"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(c *fiber.Ctx, user *model.Customer) (*model.Customer, error)
	CreateTx(tx *gorm.DB, user *model.Customer) (*model.Customer, error)
	ExistsByCif(c *fiber.Ctx, cif string) (bool, error)
	GetAssignedCustomers(marketingID uint, req *dto.AssignedCustomerRequest) ([]dto.Customer, *dto.Pagination, error)
	GetNewCustomers(req *dto.CustomerSearchRequest) ([]dto.Customer, *dto.Pagination, error)
	GetCustomerDetail(marketingID uint, customerID string) (*dto.Customer, error)
}

type customerRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCustomerRepo(db *gorm.DB, log *zap.Logger) CustomerRepository {
	return &customerRepository{
		db:  db,
		log: log,
	}
}

func (r *customerRepository) Create(c *fiber.Ctx, user *model.Customer) (*model.Customer, error) {
	error := r.db.WithContext(c.Context()).Create(user).Error
	return user, error
}

func (r *customerRepository) CreateTx(tx *gorm.DB, user *model.Customer) (*model.Customer, error) {
	if err := tx.Create(user).Error; err != nil {
		return nil, fmt.Errorf("Gagal membuat customer: %v", err)
	}

	var createdUser model.Customer
	if err := tx.Preload("CustomerProduk").Where("id = ?", user.Id).First(&createdUser).Error; err != nil {
		return nil, fmt.Errorf("Gagal mengambil data customer yang telah dibuat: %v", err)
	}
	return &createdUser, nil
}

func (r *customerRepository) ExistsByCif(c *fiber.Ctx, cif string) (bool, error) {
	var count int64
	err := r.db.WithContext(c.Context()).
		Model(&model.Customer{}).
		Where("cif = ?", cif).
		Count(&count).Error
	return count > 0, err
}

func (r *customerRepository) GetNewCustomers(req *dto.CustomerSearchRequest) ([]dto.Customer, *dto.Pagination, error) {
	var customers []dto.Customer
	var count int64

	query := r.db.Table("customers c").
		Select(`c.*`, `CASE WHEN mc.status IS NULL THEN 'new' ELSE mc.status END AS status`).
		Joins("LEFT JOIN marketing_customers mc ON mc.customer_id = c.id")

	if req.Search != "" {
		switch req.SearchBy {
		case "cif":
			query = query.Where("cif LIKE ?", "%"+req.Search+"%")
		case "nomor_hp":
			query = query.Where("nomor_hp LIKE ?", "%"+req.Search+"%")
		case "nama":
			query = query.Where("nama LIKE ?", "%"+req.Search+"%")
		case "nomor_rekening":
			query = query.Where("nomor_rekening LIKE ?", "%"+req.Search+"%")
		case "produk_eksisting":
			query = query.Where("array_to_string(produk_eksisting, ',') ILIKE ?", "%"+req.Search+"%")
		// case "product":
		// 	query = query.Where("product LIKE ?", "%"+req.Search+"%")
		default:
			query = query.Where(
				"cif ILIKE ? OR nomor_hp LIKE ? OR nama ILIKE ? OR nomor_rekening ILIKE ? OR array_to_string(produk_eksisting, ',')  ILIKE ?  OR email ILIKE ?  OR job ILIKE ?  OR address ILIKE ?",
				"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%",
				"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%",
				"%"+req.Search+"%", "%"+req.Search+"%",
			)
		}
	}

	query = query.Where("status is NULL")

	if err := query.Count(&count).Error; err != nil {
		return nil, nil, fmt.Errorf("error counting customers: %v", err)
	}

	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	if err := query.Find(&customers).Error; err != nil {
		return nil, nil, fmt.Errorf("error finding customers: %v", err)
	}

	meta := &dto.Pagination{
		CurrentPage: req.Page,
		PerPage:     req.Limit,
		TotalItems:  count,
		TotalPages:  int64(math.Ceil(float64(count) / float64(req.Limit))),
	}

	return customers, meta, nil
}

func (r *customerRepository) GetAssignedCustomers(marketingID uint, req *dto.AssignedCustomerRequest) ([]dto.Customer, *dto.Pagination, error) {
	var customers []dto.Customer
	var count int64

	query := r.db.Table("marketing_customers mc").
		Joins("JOIN customers c ON mc.customer_id = c.id").
		Where("mc.marketing_id = ?", marketingID).
		Where("mc.deleted_at IS NULL")

	if req.Search != "" {
		query = query.Where(
			"cif ILIKE ? OR nomor_hp LIKE ? OR nama ILIKE ? OR nomor_rekening ILIKE ? OR array_to_string(produk_eksisting, ',')  ILIKE ?  OR email ILIKE ?  OR job ILIKE ?  OR address ILIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%",
			"%"+req.Search+"%", "%"+req.Search+"%",
		)
	}

	if req.Status != "all" {
		query = query.Where("mc.status = ?", req.Status)
	} else {
		query = query.Where("mc.status IN ('contacted', 'rejected', 'closed')")
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, nil, fmt.Errorf("error counting assigned customers: %v", err)
	}

	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	query = query.Select(
		"c.*, mc.status, mc.notes, mc.created_at, mc.updated_at",
	)

	if err := query.Find(&customers).Error; err != nil {
		return nil, nil, fmt.Errorf("error finding assigned customers: %v", err)
	}

	meta := &dto.Pagination{
		CurrentPage: req.Page,
		PerPage:     req.Limit,
		TotalItems:  count,
		TotalPages:  int64(math.Ceil(float64(count) / float64(req.Limit))),
	}

	return customers, meta, nil
}

func (r *customerRepository) GetCustomerDetail(marketingID uint, customerID string) (*dto.Customer, error) {
	var customer dto.Customer

	var count int64
	if err := r.db.Table("customers c").
		Joins("LEFT JOIN marketing_customers mc ON c.id = mc.customer_id AND mc.marketing_id = ? AND mc.deleted_at IS NULL", marketingID).
		Where("c.cif = ? AND mc.status IS NULL", customerID).
		Or("marketing_id = ? AND c.cif = ? AND mc.deleted_at IS NULL", marketingID, customerID).
		Count(&count).Error; err != nil {
		return nil, fmt.Errorf("error checking customer: %v", err)
	}

	if count == 0 {
		return nil, errors.New("customer not found or not assigned to you")
	}

	if err := r.db.Table("customers c").
		Select(
			"c.*,"+
				"COALESCE(mc.status, 'new') as status, COALESCE(mc.notes, '') as notes, "+
				"mc.id as marketing_customer_id, mc.marketing_id, mc.product_id as product_id, "+
				"COALESCE(mc.created_at, c.created_at) as mc_created_at, mc.amount as closed_amount, "+
				"COALESCE(mc.updated_at, c.updated_at) as updated_at",
		).
		Joins("LEFT JOIN marketing_customers mc ON c.id = mc.customer_id AND mc.marketing_id = ? AND mc.deleted_at IS NULL", marketingID).
		Where("c.cif = ?", customerID).
		First(&customer).Error; err != nil {
		return nil, fmt.Errorf("error getting customer details: %v", err)
	}

	var customerProducts []struct {
		ProductID   uint   `gorm:"column:product_id"`
		ProductName string `gorm:"column:nama"`
		Icon        string `gorm:"column:ikon"`
		Prediksi    string `gorm:"column:prediksi"`
		PlafonMax   uint64 `gorm:"column:plafon_max"`
		Order       int    `gorm:"column:order"`
	}

	if err := r.db.Table("customer_products cp").
		Select("cp.product_id, p.nama, p.ikon, p.prediksi, cp.order, cp.plafon_max").
		Joins("JOIN products p ON cp.product_id = p.id").
		Where("cp.customer_id = ?", customer.Id).
		Order("cp.order ASC").
		Find(&customerProducts).Error; err != nil {
		r.log.Warn("Error fetching customer products", zap.Error(err))
	}

	product := &model.Product{}
	if err := r.db.Table("products").Where("id = ?", customer.ClosedProductID).First(&product).Error; err != nil {
		r.log.Warn("Error fetching products", zap.Error(err))
	}
	customer.ClosedProduk = *product

	customer.Produk = make([]dto.CustomerProductResponse, len(customerProducts))
	for i, cp := range customerProducts {
		customer.Produk[i] = dto.CustomerProductResponse{
			ID:        cp.ProductID,
			Nama:      cp.ProductName,
			Ikon:      cp.Icon,
			Prediksi:  cp.Prediksi,
			PlafonMax: cp.PlafonMax,
			Order:     cp.Order,
		}
	}

	return &customer, nil
}
