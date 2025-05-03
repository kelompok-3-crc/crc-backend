package repository

import (
	"fmt"
	"math"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"
	"slices"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MarketingCustomerRepository interface {
	FindByCifAndMarketingNIP(tx *gorm.DB, CIF string, NIP string) (*model.MarketingCustomer, error)
	UpdateStatusWithTx(tx *gorm.DB, mc *model.MarketingCustomer) error
	FindByCIFWithTx(tx *gorm.DB, cif string) (*model.Customer, error)
	CheckCustomerAssignmentExists(tx *gorm.DB, customerID uint64) (bool, error)
	CreateWithTx(tx *gorm.DB, mc *model.MarketingCustomer) error
	CheckAndCreateAssignment(tx *gorm.DB, customerID uint64, marketingID uint) (*model.MarketingCustomer, error)
	GetMonthlyMonitoring(month, year int) ([]dto.MarketingMonitoringResponse, error)
	// GetMarketingTargets(tx *gorm.DB, req *dto.MonitoringRequest) ([]dto.MarketingTargetDetail, error)
	GetProductPerformance(tx *gorm.DB, req *dto.ProductPerformanceRequest) (*dto.ProductPerformanceResponse, error)
}

type marketingCustomerRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewMarketingCustomerRepository(db *gorm.DB,
	log *zap.Logger) MarketingCustomerRepository {
	return &marketingCustomerRepository{db: db, log: log}
}

func (r *marketingCustomerRepository) FindByCifAndMarketingNIP(tx *gorm.DB, CIF string, NIP string) (*model.MarketingCustomer, error) {
	var mc model.MarketingCustomer
	err := tx.Table("marketing_customers").
		Joins("JOIN customers ON customers.id = marketing_customers.customer_id").
		Joins("JOIN users ON users.id = marketing_customers.marketing_id").
		Where("customers.cif = ? AND users.nip = ? AND marketing_customers.deleted_at IS NULL",
			CIF, NIP).
		First(&mc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("error finding marketing customer: %v", err)
	}
	return &mc, nil
}

func (r *marketingCustomerRepository) UpdateStatusWithTx(tx *gorm.DB, mc *model.MarketingCustomer) error {
	return tx.Save(mc).Error
}

func (r *marketingCustomerRepository) FindByCIFWithTx(tx *gorm.DB, cif string) (*model.Customer, error) {
	var customer model.Customer
	err := tx.Where("cif = ?", cif).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *marketingCustomerRepository) CheckCustomerAssignmentExists(tx *gorm.DB, customerID uint64) (bool, error) {
	var count int64
	err := tx.Model(&model.MarketingCustomer{}).
		Where("customer_id = ? AND deleted_at IS NULL", customerID).
		Count(&count).Error
	return count > 0, err
}

func (r *marketingCustomerRepository) CreateWithTx(tx *gorm.DB, mc *model.MarketingCustomer) error {
	return tx.Create(mc).Error
}
func (r *marketingCustomerRepository) CheckAndCreateAssignment(tx *gorm.DB, customerID uint64, marketingID uint) (*model.MarketingCustomer, error) {
	// Use a single query with FOR UPDATE to prevent race conditions
	var count int64
	err := tx.Model(&model.MarketingCustomer{}).
		Where("customer_id = ? AND deleted_at IS NULL", customerID).
		Count(&count).Error

	if err != nil {
		return nil, fmt.Errorf("Gagal mengecek assignment : %v", err)
	}

	if count > 0 {
		return nil, fmt.Errorf("Customer telah memiliki assignment lain")
	}

	// Create new assignment since none exists
	newAssignment := &model.MarketingCustomer{
		CustomerID:  customerID,
		MarketingID: marketingID,
		Status:      string(model.CustomerStatusNew),
	}

	if err := tx.Create(newAssignment).Error; err != nil {
		return nil, fmt.Errorf("Gagal membuat assignment: %v", err)
	}

	return newAssignment, nil
}

// func (r *marketingCustomerRepository) GetMonthlyMonitoring(tx *gorm.DB, month, year int) ([]dto.MarketingMonitoringResponse, error) {
// 	var result []dto.MarketingMonitoringResponse

// 	err := tx.Raw(`
//         WITH monthly_closings AS (
//             SELECT
//                 u.id,
//                 COALESCE(SUM(mc.amount), 0) as closing_amount
//             FROM users u
//             LEFT JOIN marketing_customers mc ON mc.marketing_id = u.id
//                 AND mc.status = 'closed'
//                 AND mc.deleted_at IS NULL
//                 AND EXTRACT(MONTH FROM mc.updated_at) = ?
//                 AND EXTRACT(YEAR FROM mc.updated_at) = ?
//             WHERE u.role = 'marketing'
//             AND u.deleted_at IS NULL
//             GROUP BY u.id
//         ),
//         monthly_targets AS (
//             SELECT
//                 u.id,
//                 COALESCE(SUM(mt.target_amount), 0) as target_amount
//             FROM users u
//             LEFT JOIN marketing_target_bulanan mt ON mt.marketing_id = u.id
//                 AND mt.bulan = ?
//                 AND mt.tahun = ?
//                 AND mt.deleted_at IS NULL
//             WHERE u.role = 'marketing'
//             AND u.deleted_at IS NULL
//             GROUP BY u.id
//         ),
//         previous_excess AS (
//             SELECT
//                 u.id,
//                 GREATEST(
//                     COALESCE(SUM(mc.amount), 0) -
//                     COALESCE((
//                         SELECT SUM(mt.target_amount)
//                         FROM marketing_target_bulanan mt
//                         WHERE mt.marketing_id = u.id
//                         AND mt.bulan = ?
//                         AND mt.tahun = ?
//                         AND mt.deleted_at IS NULL
//                     ), 0),
//                     0
//                 ) as carry_over
//             FROM users u
//             LEFT JOIN marketing_customers mc ON mc.marketing_id = u.id
//                 AND mc.status = 'closed'
//                 AND mc.deleted_at IS NULL
//                 AND EXTRACT(MONTH FROM mc.updated_at) = ?
//                 AND EXTRACT(YEAR FROM mc.updated_at) = ?
//             WHERE u.role = 'marketing'
//             AND u.deleted_at IS NULL
//             GROUP BY u.id
//         )
//         SELECT
//             u.nip as marketing_nip,
//             u.nama as marketing_name,
//             COALESCE(mt.target_amount, 0) as monthly_target,
//             COALESCE(mc.closing_amount, 0) as monthly_closing,
//             COALESCE(pe.carry_over, 0) as carry_over,
//             CASE
//                 WHEN COALESCE(mt.target_amount, 0) > 0
//                 THEN ROUND(((COALESCE(mc.closing_amount, 0) + COALESCE(pe.carry_over, 0)) / mt.target_amount * 100), 2)
//                 ELSE 0
//             END as achievement_percentage
//         FROM users u
//         LEFT JOIN monthly_closings mc ON mc.id = u.id
//         LEFT JOIN monthly_targets mt ON mt.id = u.id
//         LEFT JOIN previous_excess pe ON pe.id = u.id
//         WHERE u.role = 'marketing'
//         AND u.deleted_at IS NULL
//         ORDER BY u.nama
//     `, month, year, month, year, month-1, year, month-1, year).
// 		Scan(&result).Error

// 	if err != nil {
// 		return nil, fmt.Errorf("error getting monitoring data: %v", err)
// 	}

// 	return result, nil
// }

func (r *marketingCustomerRepository) GetMonthlyMonitoring(month, year int) ([]dto.MarketingMonitoringResponse, error) {
	var tempResult []struct {
		MarketingNIP   string  `gorm:"column:marketing_nip"`
		MarketingName  string  `gorm:"column:marketing_name"`
		MonthlyTarget  float64 `gorm:"column:monthly_target"`
		MonthlyClosing float64 `gorm:"column:monthly_closing"`
		CarryOver      float64 `gorm:"column:carry_over"`
	}

	// Main query remains the same until the SELECT part
	err := r.db.Raw(`
      WITH monthly_closings AS (
        SELECT 
            u.id,
            COALESCE(SUM(mc.amount), 0) as closing_amount
        FROM users u
        LEFT JOIN marketing_customers mc ON mc.marketing_id = u.id 
            AND mc.status = 'closed'
            AND mc.deleted_at IS NULL
            AND EXTRACT(MONTH FROM mc.updated_at) = ?
            AND EXTRACT(YEAR FROM mc.updated_at) = ?
        WHERE u.role = 'marketing'
        AND u.deleted_at IS NULL
        GROUP BY u.id
    ),
    monthly_targets AS (
        SELECT 
            u.id,
            COALESCE(SUM(mt.target_amount), 0) as target_amount
        FROM users u
        LEFT JOIN marketing_target_bulanan mt ON mt.marketing_id = u.id
            AND mt.bulan = ?
            AND mt.tahun = ?
            AND mt.deleted_at IS NULL
        WHERE u.role = 'marketing'
        AND u.deleted_at IS NULL
        GROUP BY u.id
    ),
    previous_excess AS (
        SELECT 
            u.id,
            GREATEST(
                COALESCE(SUM(mc.amount), 0) - 
                COALESCE((
                    SELECT SUM(mt.target_amount) 
                    FROM marketing_target_bulanan mt 
                    WHERE mt.marketing_id = u.id 
                    AND mt.bulan = ? 
                    AND mt.tahun = ?
                    AND mt.deleted_at IS NULL
                ), 0),
                0
            ) as carry_over
        FROM users u
        LEFT JOIN marketing_customers mc ON mc.marketing_id = u.id
            AND mc.status = 'closed'
            AND mc.deleted_at IS NULL
            AND EXTRACT(MONTH FROM mc.updated_at) = ?
            AND EXTRACT(YEAR FROM mc.updated_at) = ?
        WHERE u.role = 'marketing'
        AND u.deleted_at IS NULL
        GROUP BY u.id
    )
        SELECT 
            u.nip as marketing_nip,
            u.nama as marketing_name,
            COALESCE(mt.target_amount, 0) as monthly_target,
            COALESCE(mc.closing_amount, 0) as monthly_closing,
            COALESCE(pe.carry_over, 0) as carry_over
        FROM users u
        LEFT JOIN monthly_closings mc ON mc.id = u.id
        LEFT JOIN monthly_targets mt ON mt.id = u.id
        LEFT JOIN previous_excess pe ON pe.id = u.id
        WHERE u.role = 'marketing'
        AND u.deleted_at IS NULL
        ORDER BY monthly_closing DESC
    `, month, year, month, year, month-1, year, month-1, year).
		Scan(&tempResult).Error

	if err != nil {
		return nil, fmt.Errorf("error getting monitoring data: %v", err)
	}

	// Convert to ChartJS format
	result := make([]dto.MarketingMonitoringResponse, len(tempResult))
	for i, temp := range tempResult {
		monthlyAchieved := temp.MonthlyClosing + temp.CarryOver

		// Main chart data
		result[i] = dto.MarketingMonitoringResponse{
			MarketingNIP:    temp.MarketingNIP,
			MarketingName:   temp.MarketingName,
			MonthlyAchieved: monthlyAchieved,
			MonthlyTarget:   temp.MonthlyTarget,
			Labels:          []string{"Tercapai", "Belum"},
			Datasets: []dto.ChartDataset{
				{
					Data:            []float64{monthlyAchieved, math.Max(temp.MonthlyTarget-monthlyAchieved, 0)},
					BackgroundColor: []string{"#14b8a6", "#e5e7eb"},
					BorderWidth:     0,
				},
			},
		}

		var productDetails []struct {
			ProductID       uint    `gorm:"column:product_id"`
			ProductName     string  `gorm:"column:product_name"`
			MonthlyTarget   float64 `gorm:"column:monthly_target"`
			MonthlyAchieved float64 `gorm:"column:monthly_achieved"`
		}

		err := r.db.Raw(`
			WITH product_data AS (
				SELECT 
					p.id as product_id,
					p.nama as product_name,
					COALESCE(mt.target_amount, 0) as monthly_target,
					(COALESCE(SUM(mc.amount), 0) + 
					COALESCE(
						GREATEST(
							COALESCE((
								SELECT SUM(pmc.amount)
								FROM marketing_customers pmc
								WHERE pmc.marketing_id = (SELECT id FROM users WHERE nip = ? AND deleted_at IS NULL)
								AND pmc.product_id = p.id
								AND pmc.status = 'closed'
								AND EXTRACT(MONTH FROM pmc.updated_at) = ?
								AND EXTRACT(YEAR FROM pmc.updated_at) = ?
								AND pmc.deleted_at IS NULL
							), 0) - 
							COALESCE((
								SELECT target_amount
								FROM marketing_target_bulanan
								WHERE marketing_id = (SELECT id FROM users WHERE nip = ? AND deleted_at IS NULL)
								AND product_id = p.id
								AND bulan = ?
								AND tahun = ?
								AND deleted_at IS NULL
							), 0),
							0
						),
						0
					)) as monthly_achieved
				FROM products p
				LEFT JOIN marketing_target_bulanan mt ON mt.product_id = p.id
					AND mt.marketing_id = (SELECT id FROM users WHERE nip = ? AND deleted_at IS NULL)
					AND mt.bulan = ? 
					AND mt.tahun = ?
					AND mt.deleted_at IS NULL
				LEFT JOIN marketing_customers mc ON mc.product_id = p.id
					AND mc.marketing_id = mt.marketing_id
					AND mc.status = 'closed'
					AND EXTRACT(MONTH FROM mc.updated_at) = ?
					AND EXTRACT(YEAR FROM mc.updated_at) = ?
					AND mc.deleted_at IS NULL
				WHERE p.deleted_at IS NULL
				GROUP BY p.id, p.nama, mt.target_amount
				ORDER BY p.id
			)
			SELECT * FROM product_data
		`, temp.MarketingNIP, month-1, year,
			temp.MarketingNIP, month-1, year,
			temp.MarketingNIP, month, year,
			month, year).Scan(&productDetails).Error

		if err != nil {
			return nil, fmt.Errorf("error getting product details: %v", err)
		}

		// Convert product details to chart format
		result[i].TargetDetails = make([]dto.ProductChart, len(productDetails))
		for j, prod := range productDetails {
			result[i].TargetDetails[j] = dto.ProductChart{
				ProductID:       prod.ProductID,
				ProductName:     prod.ProductName,
				MonthlyAchieved: prod.MonthlyAchieved,
				MonthlyTarget:   prod.MonthlyTarget,
				Labels:          []string{"Tercapai", "Belum"},
				Datasets: []dto.ChartDataset{
					{
						Data:            []float64{prod.MonthlyAchieved, math.Max(prod.MonthlyTarget-prod.MonthlyAchieved, 0)},
						BackgroundColor: []string{"#14b8a6", "#e5e7eb"},
						BorderWidth:     0,
					},
				},
			}
		}
	}

	return result, nil
}

func (r *marketingCustomerRepository) GetProductPerformance(tx *gorm.DB, req *dto.ProductPerformanceRequest) (*dto.ProductPerformanceResponse, error) {
	// Define time format based on grouping option
	timeFormat := ""
	switch req.GroupBy {
	case "week":
		timeFormat = "YYYY-WW"
	case "month":
		timeFormat = "YYYY-MM"
	case "year":
		timeFormat = "YYYY"
	}

	// Get product info
	var product struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:nama"`
	}

	err := tx.Table("products").
		Select("id, nama").
		Where("id = ? AND deleted_at IS NULL", req.ProductID).
		First(&product).Error

	if err != nil {
		return nil, fmt.Errorf("error getting product data: %v", err)
	}

	// Query for time periods and marketing achievements
	var timeSeries []struct {
		TimeLabel     string  `gorm:"column:time_label"`
		MarketingNIP  string  `gorm:"column:marketing_nip"`
		MarketingName string  `gorm:"column:marketing_name"`
		Achievement   float64 `gorm:"column:achievement"`
	}

	// Query with branch_total (sum of all marketing achievements)
	err = tx.Raw(`
    WITH RECURSIVE time_periods AS (
        SELECT 
            TO_CHAR(d, ?) as time_label
        FROM generate_series(
            ?::timestamp, 
            ?::timestamp, 
            CASE ? 
                WHEN 'week' THEN '1 week'::interval
                WHEN 'month' THEN '1 month'::interval
                ELSE '1 year'::interval
            END
        ) d
    ),
    marketing_achievements AS (
        SELECT 
            tp.time_label,
            u.nip as marketing_nip,
            u.nama as marketing_name,
            COALESCE(SUM(mc.amount), 0) as achievement
        FROM time_periods tp
        CROSS JOIN users u
        LEFT JOIN marketing_customers mc ON mc.marketing_id = u.id
            AND mc.product_id = ?
            AND mc.status = 'closed'
            AND TO_CHAR(mc.updated_at, ?) = tp.time_label
            AND mc.deleted_at IS NULL
        WHERE u.role = 'marketing'
        AND u.deleted_at IS NULL
        GROUP BY tp.time_label, u.id, u.nip, u.nama
    ),
    branch_total AS (
        SELECT 
            tp.time_label,
            'ALL' as marketing_nip,
            'All' as marketing_name,
            COALESCE(SUM(ma.achievement), 0) as achievement
        FROM time_periods tp
        LEFT JOIN marketing_achievements ma ON ma.time_label = tp.time_label
        GROUP BY tp.time_label
    )
    SELECT * FROM branch_total
    UNION ALL
    SELECT * FROM marketing_achievements
    ORDER BY time_label, marketing_nip
    `, timeFormat, req.StartDate, req.EndDate, req.GroupBy,
		req.ProductID, timeFormat).Scan(&timeSeries).Error

	if err != nil {
		return nil, fmt.Errorf("error getting time series data: %v", err)
	}

	// Organize data for ChartJS
	timeLabels := make([]string, 0)
	marketingMap := make(map[string]*dto.MarketingDataset)
	colors := []string{"#0ea5e9", "#14b8a6", "#f59e0b", "#8b5cf6", "#ec4899", "#ef4444"}

	// First pass: collect unique labels and initialize datasets
	for _, ts := range timeSeries {
		if !slices.Contains(timeLabels, ts.TimeLabel) {
			timeLabels = append(timeLabels, ts.TimeLabel)
		}

		if _, exists := marketingMap[ts.MarketingNIP]; !exists {
			colorIndex := len(marketingMap) % len(colors)
			marketingMap[ts.MarketingNIP] = &dto.MarketingDataset{
				MarketingNIP:    ts.MarketingNIP,
				Label:           ts.MarketingName,
				Data:            make([]float64, 0, len(timeLabels)),
				BorderColor:     colors[colorIndex],
				BackgroundColor: colors[colorIndex],
				Fill:            false,
				Tension:         0.1,
			}
		}
	}

	// Second pass: fill in data points
	for _, label := range timeLabels {
		for nip, dataset := range marketingMap {
			found := false
			for _, ts := range timeSeries {
				if ts.TimeLabel == label && ts.MarketingNIP == nip {
					dataset.Data = append(dataset.Data, ts.Achievement)
					found = true
					break
				}
			}
			if !found {
				dataset.Data = append(dataset.Data, 0)
			}
		}
	}

	// When converting to the final slice, ensure "ALL" appears first
	datasets := make([]dto.MarketingDataset, 0, len(marketingMap))

	// Add ALL dataset first if it exists
	if allDataset, exists := marketingMap["ALL"]; exists {
		datasets = append(datasets, *allDataset)
		delete(marketingMap, "ALL") // Remove from map to avoid adding it twice
	}

	// Then add all other datasets
	for _, dataset := range marketingMap {
		datasets = append(datasets, *dataset)
	}

	return &dto.ProductPerformanceResponse{
		Produk:   product.Name,
		ProdukID: strconv.FormatUint(uint64(product.ID), 10),
		Labels:   timeLabels,
		Datasets: datasets,
	}, nil
}
