package usecase

import (
	"context"
	"fmt"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"
	"ml-prediction/internal/app/repository"

	"gorm.io/gorm"
)

type MarketingCustomerUsecase interface {
	UpdateCustomerStatus(req *dto.UpdateCustomerStatusRequest, marketingNIP string) error
	GetMonthlyMonitoring(req *dto.MonitoringRequest) ([]dto.MarketingMonitoringResponse, error)
	GetMonthlyMonitoringMarketing(req *dto.MonitoringRequest) (*dto.MarketingMonitoringResponse, error)
	GetProductPerformance(ctx context.Context, req *dto.ProductPerformanceRequest) (*dto.ProductPerformanceResponse, error)
}

type marketingCustomerUsecase struct {
	marketingCustomerRepo repository.MarketingCustomerRepository
	userRepo              repository.UserRepository
	db                    *gorm.DB
}

func NewMarketingCustomerUsecase(
	mcRepo repository.MarketingCustomerRepository,
	userRepo repository.UserRepository,
	db *gorm.DB,
) MarketingCustomerUsecase {
	return &marketingCustomerUsecase{
		marketingCustomerRepo: mcRepo,
		userRepo:              userRepo,
		db:                    db,
	}
}

func (u *marketingCustomerUsecase) UpdateCustomerStatus(req *dto.UpdateCustomerStatusRequest, marketingNIP string) error {
	tx := u.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get marketing user
	marketing, err := u.userRepo.FindByNIPWithTx(tx, marketingNIP)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("marketing tidak ditemukan: %v", err)
	}

	// Find customer
	customer, err := u.marketingCustomerRepo.FindByCIFWithTx(tx, req.CIF)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("customer dengan CIF %s tidak ditemukan: %v", req.CIF, err)
	}

	// Try to find existing assignment
	mc, err := u.marketingCustomerRepo.FindByCifAndMarketingNIP(tx, customer.CIF, marketing.NIP)
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return fmt.Errorf("gagal mencari assignment: %v", err)
	}

	// If no assignment exists, check if customer can be assigned
	if err == gorm.ErrRecordNotFound {
		mc, err = u.marketingCustomerRepo.CheckAndCreateAssignment(tx, customer.Id, marketing.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal membuat assignment: %v", err)
		}
	}
	// Check if status is final
	if mc.Status == string(model.CustomerStatusClosed) ||
		mc.Status == string(model.CustomerStatusRejected) {
		tx.Rollback()
		return fmt.Errorf("tidak dapat mengubah status yang sudah final")
	}

	// Validate status requirements
	switch req.Status {
	case string(model.CustomerStatusClosed):
		if req.ProductID == nil || req.Amount == nil {
			tx.Rollback()
			return fmt.Errorf("product dan amount harus diisi untuk status closed")
		}
	case string(model.CustomerStatusRejected):
		if req.Notes == "" {
			tx.Rollback()
			return fmt.Errorf("notes harus diisi untuk status rejected")
		}
	}

	// Update status and related fields
	mc.Status = req.Status
	if req.Status == string(model.CustomerStatusClosed) {
		mc.ProductID = req.ProductID
		mc.Amount = req.Amount
	} else if req.Status == string(model.CustomerStatusRejected) {
		mc.Notes = req.Notes
	}

	if err := u.marketingCustomerRepo.UpdateStatusWithTx(tx, mc); err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mengupdate status: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menyimpan target: %v", err)
	}
	return nil
}

func (u *marketingCustomerUsecase) GetMonthlyMonitoring(req *dto.MonitoringRequest) ([]dto.MarketingMonitoringResponse, error) {

	result, err := u.marketingCustomerRepo.GetMonthlyMonitoring(req.Month, req.Year)
	return result, err
}

func (u *marketingCustomerUsecase) GetMonthlyMonitoringMarketing(req *dto.MonitoringRequest) (*dto.MarketingMonitoringResponse, error) {
	user, err := u.userRepo.FindByNIPWithTx(u.db, req.NIP)
	if err != nil || user.Role != "marketing" {
		return nil, fmt.Errorf("user tidak ditemukan atau bukan marketing: %v", err)
	}

	result, err := u.marketingCustomerRepo.GetMonthlyMonitoring(req.Month, req.Year)
	return &result[0], err
}

func (s *marketingCustomerUsecase) GetProductPerformance(ctx context.Context, req *dto.ProductPerformanceRequest) (*dto.ProductPerformanceResponse, error) {
	return s.marketingCustomerRepo.GetProductPerformance(s.db, req)
}
