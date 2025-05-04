package repository

import (
	"fmt"
	"math"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TargetRepository interface {
	IsTargetTahunanExist(tahun int, kantorID, productID uint) (bool, error)
	CreateTargetTahunan(t *model.TargetProdukTahunan) error
	CreateTargetBulanan(b *model.TargetProdukBulanan) error
	GetTargetBulanan(tahun int, bulan int, kantorID, productID uint) (*model.TargetProdukBulanan, error)
	GetMarketingTargetSum(tahun int, bulan int, kantorID, productID uint) (int64, error)
	CreateMarketingTarget(t *model.MarketingTargetBulanan) error
	ListMarketingInKantorCabang(kantorID uint) ([]model.User, error)
	GetTargetBulananWithTx(tx *gorm.DB, tahun, bulan int, kantorCabangID, productID uint) (*model.TargetProdukBulanan, error)
	GetMarketingTargetSumWithTx(tx *gorm.DB, tahun, bulan int, kantorCabangID, productID uint, excludeMarketingID uint) (int64, error)
	CreateMarketingTargetWithTx(tx *gorm.DB, target *model.MarketingTargetBulanan) error
	UpdateMarketingTargetWithTx(tx *gorm.DB, target *model.MarketingTargetBulanan) error
	GetExistingMarketingTargetsWithTx(tx *gorm.DB, tahun, bulan int, marketingID uint) ([]model.MarketingTargetBulanan, error)
	GetExistingTargetTahunanWithTx(tx *gorm.DB, tahun int, kantorCabangID uint) ([]model.TargetProdukTahunan, error)
	UpdateTargetTahunanWithTx(tx *gorm.DB, target *model.TargetProdukTahunan) error
	CreateTargetTahunanWithTx(tx *gorm.DB, target *model.TargetProdukTahunan) error
	UpdateTargetBulananWithTx(tx *gorm.DB, target *model.TargetProdukBulanan) error
	CreateTargetBulananWithTx(tx *gorm.DB, target *model.TargetProdukBulanan) error
	GetMarketingTargets(req *dto.MonitoringRequest) ([]dto.MarketingTargetDetail, error)
	GetTargetSummary(tx *gorm.DB, userID uint, role string, month, year int) (*dto.TargetSummaryResponse, error)
}

type targetRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewTargetRepository(db *gorm.DB, log *zap.Logger) TargetRepository {
	return &targetRepository{db: db, log: log}
}

func (r *targetRepository) IsTargetTahunanExist(tahun int, kantorID, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.TargetProdukTahunan{}).
		Where("tahun = ? AND kantor_cabang_id = ? AND product_id = ?", tahun, kantorID, productID).
		Count(&count).Error
	return count > 0, err
}

func (r *targetRepository) CreateTargetTahunan(t *model.TargetProdukTahunan) error {
	return r.db.Create(t).Error
}

func (r *targetRepository) CreateTargetBulanan(b *model.TargetProdukBulanan) error {
	return r.db.Create(b).Error
}

func (r *targetRepository) GetTargetBulanan(tahun int, bulan int, kantorID, productID uint) (*model.TargetProdukBulanan, error) {
	var target model.TargetProdukBulanan
	err := r.db.Where("tahun = ? AND bulan = ? AND kantor_cabang_id = ? AND product_id = ?",
		tahun, bulan, kantorID, productID).First(&target).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *targetRepository) GetMarketingTargetSum(tahun int, bulan int, kantorID, productID uint) (int64, error) {
	var sum int64
	err := r.db.Model(&model.MarketingTargetBulanan{}).
		Where("tahun = ? AND bulan = ? AND kantor_cabang_id = ? AND product_id = ?",
			tahun, bulan, kantorID, productID).
		Select("COALESCE(SUM(target_amount), 0)").
		Scan(&sum).Error
	return sum, err
}

func (r *targetRepository) CreateMarketingTarget(t *model.MarketingTargetBulanan) error {
	return r.db.Create(t).Error
}

func (r *targetRepository) ListMarketingInKantorCabang(kantorID uint) ([]model.User, error) {
	var marketings []model.User
	err := r.db.Where("role = ? AND kantor_cabang_id = ?", "marketing", kantorID).
		Find(&marketings).Error
	return marketings, err
}

func (r *targetRepository) GetTargetBulananWithTx(tx *gorm.DB, tahun, bulan int, kantorCabangID, productID uint) (*model.TargetProdukBulanan, error) {
	var target model.TargetProdukBulanan
	err := tx.Where("tahun = ? AND bulan = ? AND kantor_cabang_id = ? AND product_id = ?",
		tahun, bulan, kantorCabangID, productID).First(&target).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *targetRepository) GetMarketingTargetSumWithTx(tx *gorm.DB, tahun, bulan int, kantorCabangID, productID uint, excludeMarketingID uint) (int64, error) {
	var sum int64
	err := tx.Model(&model.MarketingTargetBulanan{}).
		Where("tahun = ? AND bulan = ? AND kantor_cabang_id = ? AND product_id = ? AND marketing_id <> ?",
			tahun, bulan, kantorCabangID, productID, excludeMarketingID).
		Select("COALESCE(SUM(target_amount), 0)").
		Scan(&sum).Error
	return sum, err
}
func (r *targetRepository) CreateMarketingTargetWithTx(tx *gorm.DB, target *model.MarketingTargetBulanan) error {
	return tx.Create(target).Error
}
func (r *targetRepository) UpdateMarketingTargetWithTx(tx *gorm.DB, target *model.MarketingTargetBulanan) error {
	return tx.Save(target).Error
}

func (r *targetRepository) GetExistingMarketingTargetsWithTx(tx *gorm.DB, tahun, bulan int, marketingID uint) ([]model.MarketingTargetBulanan, error) {
	var targets []model.MarketingTargetBulanan
	err := tx.Where("tahun = ? AND bulan = ? AND marketing_id = ?",
		tahun, bulan, marketingID).Find(&targets).Error
	return targets, err
}

func (r *targetRepository) GetExistingTargetTahunanWithTx(tx *gorm.DB, tahun int, kantorCabangID uint) ([]model.TargetProdukTahunan, error) {
	var targets []model.TargetProdukTahunan
	err := tx.Where("tahun = ? AND kantor_cabang_id = ?",
		tahun, kantorCabangID).Find(&targets).Error
	return targets, err
}

func (r *targetRepository) UpdateTargetTahunanWithTx(tx *gorm.DB, target *model.TargetProdukTahunan) error {
	return tx.Save(target).Error
}

func (r *targetRepository) CreateTargetTahunanWithTx(tx *gorm.DB, target *model.TargetProdukTahunan) error {
	return tx.Create(target).Error
}

func (r *targetRepository) UpdateTargetBulananWithTx(tx *gorm.DB, target *model.TargetProdukBulanan) error {
	return tx.Save(target).Error
}

func (r *targetRepository) CreateTargetBulananWithTx(tx *gorm.DB, target *model.TargetProdukBulanan) error {
	return tx.Create(target).Error
}

func (r *targetRepository) GetMarketingTargets(req *dto.MonitoringRequest) ([]dto.MarketingTargetDetail, error) {
	type marketingTargetTemp struct {
		MarketingNIP  string  `gorm:"column:marketing_nip"`
		MarketingName string  `gorm:"column:marketing_name"`
		HasTarget     bool    `gorm:"column:has_target"`
		TotalTarget   float64 `gorm:"column:total_target"`
	}

	var tempDetails []marketingTargetTemp

	query := `
        WITH marketing_targets AS (
            SELECT 
                u.id,
                u.nip as marketing_nip,
                u.nama as marketing_name,
                CASE WHEN SUM(mt.target_amount) > 0 THEN true ELSE false END as has_target,
                COALESCE(SUM(mt.target_amount), 0) as total_target
            FROM users u
            LEFT JOIN marketing_target_bulanan mt ON mt.marketing_id = u.id 
                AND mt.bulan = ? 
                AND mt.tahun = ?
                AND mt.deleted_at IS NULL
            WHERE u.role = 'marketing'
                AND u.deleted_at IS NULL
            GROUP BY u.id, u.nip, u.nama
        )
        SELECT * FROM marketing_targets mt
        WHERE 1=1
    `
	args := []interface{}{req.Month, req.Year}

	if req.HasTarget != nil {
		query += " AND mt.has_target = ?"
		args = append(args, *req.HasTarget)
	}

	if req.Search != "" {
		query += " AND (mt.marketing_nip ILIKE ? OR mt.marketing_name ILIKE ?)"
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if req.NIP != "" {
		query += " AND (mt.marketing_nip = ?)"
		args = append(args, req.NIP)
	}

	query += " ORDER BY mt.marketing_name"

	err := r.db.Raw(query, args...).Scan(&tempDetails).Error
	if err != nil {
		return nil, fmt.Errorf("error getting target data: %v", err)
	}

	// Convert to final response structure
	details := make([]dto.MarketingTargetDetail, len(tempDetails))
	for i, temp := range tempDetails {
		details[i] = dto.MarketingTargetDetail{
			MarketingNIP:  temp.MarketingNIP,
			MarketingName: temp.MarketingName,
			HasTarget:     temp.HasTarget,
			TotalTarget:   temp.TotalTarget,
		}

		var productTargets []dto.ProductTarget
		err := r.db.Raw(`
				SELECT 
					p.id as product_id,
					p.nama as product_name,
					COALESCE(mt.target_amount, 0) as amount
				FROM products p
				LEFT JOIN (
					SELECT product_id, target_amount 
					FROM marketing_target_bulanan 
					WHERE marketing_id = (
						SELECT id FROM users WHERE nip = ? AND deleted_at IS NULL
					)
					AND bulan = ?
					AND tahun = ?
					AND deleted_at IS NULL
				) mt ON mt.product_id = p.id
				WHERE p.deleted_at IS NULL
				ORDER BY p.id ASC
			`, temp.MarketingNIP, req.Month, req.Year).Scan(&productTargets).Error

		if err != nil {
			return nil, fmt.Errorf("error getting product targets: %v", err)
		}
		details[i].TargetDetails = productTargets
	}

	return details, nil
}

func (r *targetRepository) GetTargetSummary(tx *gorm.DB, userID uint, role string, month, year int) (*dto.TargetSummaryResponse, error) {
	var response dto.TargetSummaryResponse

	// Get user basic info (used for both roles)
	var user struct {
		ID       uint   `gorm:"column:id"`
		Name     string `gorm:"column:nama"`
		NIP      string `gorm:"column:nip"`
		BranchID uint   `gorm:"column:kantor_cabang_id"`
	}

	if err := tx.Table("users").
		Select("id, nama, nip, kantor_cabang_id").
		Where("id = ? AND deleted_at IS NULL", userID).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("error getting user data: %v", err)
	}

	var query string
	var args []interface{}

	if role == "bm" {
		query = `
			WITH branch_achievements AS (
				SELECT 
					mc.product_id,
					SUM(mc.amount) as amount
				FROM marketing_customers mc
				JOIN users u ON mc.marketing_id = u.id AND u.kantor_cabang_id = ?
				WHERE mc.status = 'closed'
				AND EXTRACT(MONTH FROM mc.updated_at) = ?
				AND EXTRACT(YEAR FROM mc.updated_at) = ?
				AND mc.deleted_at IS NULL
				GROUP BY mc.product_id
			)
			SELECT 
				COALESCE(SUM(tb.target_amount), 0) as target,
				COALESCE((
					SELECT SUM(ba.amount)
					FROM branch_achievements ba
					WHERE ba.product_id IN (
						SELECT product_id FROM target_produk_bulanan 
						WHERE bulan = ? AND tahun = ? AND kantor_cabang_id = ?
					)
				), 0) as achieved
			FROM target_produk_bulanan tb
			WHERE tb.bulan = ? AND tb.tahun = ? AND tb.kantor_cabang_id = ?
			AND tb.deleted_at IS NULL
			`
		args = []interface{}{user.BranchID, month, year, month, year, user.BranchID, month, year, user.BranchID}
	} else {
		query = `
				SELECT 
					COALESCE(SUM(mt.target_amount), 0) as target,
					COALESCE(SUM(CASE WHEN mc.status = 'closed' THEN mc.amount ELSE 0 END), 0) as achieved
				FROM (
					SELECT ?::integer as marketing_id, ?::integer as month, ?::integer as year
				) params
				LEFT JOIN marketing_target_bulanan mt ON mt.marketing_id = params.marketing_id
					AND mt.bulan = params.month AND mt.tahun = params.year
					AND mt.deleted_at IS NULL
				LEFT JOIN marketing_customers mc ON mc.marketing_id = params.marketing_id
					AND mc.status = 'closed'
					AND EXTRACT(MONTH FROM mc.updated_at) = params.month
					AND EXTRACT(YEAR FROM mc.updated_at) = params.year
					AND mc.deleted_at IS NULL
			`
		args = []interface{}{userID, month, year}
	}

	var products []struct {
		ProductID   uint    `gorm:"column:product_id"`
		ProductName string  `gorm:"column:product_name"`
		Target      float64 `gorm:"column:target"`
		Achieved    float64 `gorm:"column:achieved"`
	}

	if err := tx.Raw(query, args...).Scan(&products).Error; err != nil {
		return nil, fmt.Errorf("error getting product targets: %v", err)
	}

	for _, p := range products {
		response.TotalTarget += p.Target
		response.Achieved += p.Achieved
	}

	var branch struct {
		Nama string `gorm:"column:nama"`
	}
	if err := r.db.Table("kantor_cabang").
		Where("id = ?", user.BranchID).
		First(&branch).Error; err != nil {
		r.log.Warn("Failed to get branch name", zap.Error(err), zap.Uint("branchID", user.BranchID))
		response.BranchName = "Unknown Branch"
	} else {
		response.BranchName = branch.Nama
	}
	response.Type = "bm"
	if role == "marketing" {
		response.Type = "marketing"
	}

	response.TargetMonth = month
	response.TargetYear = year
	response.TargetSetted = response.TotalTarget > 0
	response.Name = user.Name
	response.NIP = user.NIP

	response.Percentage = calculatePercentage(response.Achieved, response.TotalTarget)

	return &response, nil
}

// Helper function to calculate percentage safely
func calculatePercentage(achieved, target float64) float64 {
	if target <= 0 {
		return 0
	}
	return math.Round((achieved/target)*100*100) / 100 // Round to 2 decimal places
}
