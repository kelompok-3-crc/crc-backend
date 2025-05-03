package usecase

import (
	"fmt"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"
	"ml-prediction/internal/app/repository"

	"gorm.io/gorm"
)

type MarketingTargetUsecase interface {
	AssignBulkMarketingTarget(req *dto.AssignMarketingTargetRequest, userNIP string) error
	GetMarketingTargets(req *dto.MonitoringRequest) ([]dto.MarketingTargetDetail, error)
}

type marketingTargetUsecase struct {
	targetRepo repository.TargetRepository
	userRepo   repository.UserRepository
	db         *gorm.DB
}

func NewMarketingTargetUsecase(targetRepo repository.TargetRepository, userRepo repository.UserRepository, db *gorm.DB) MarketingTargetUsecase {
	return &marketingTargetUsecase{
		targetRepo: targetRepo,
		userRepo:   userRepo,
		db:         db,
	}
}
func (u *marketingTargetUsecase) AssignBulkMarketingTarget(req *dto.AssignMarketingTargetRequest, userNIP string) error {
	tx := u.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	bm, err := u.userRepo.FindByNIPWithTx(tx, userNIP)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("user tidak ditemukan")
	}
	if bm.Role != "bm" {
		tx.Rollback()
		return fmt.Errorf("hanya Branch Manager yang dapat mengatur target marketing")
	}

	marketing, err := u.userRepo.FindByNIPWithTx(tx, req.MarketingNIP)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("marketing tidak ditemukan")
	}
	if marketing.KantorCabangID == nil || *marketing.KantorCabangID != *bm.KantorCabangID {
		tx.Rollback()
		return fmt.Errorf("marketing tidak berada dalam kantor cabang yang sama")
	}

	existingTargets, err := u.targetRepo.GetExistingMarketingTargetsWithTx(tx, req.Tahun, req.Bulan, marketing.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal memeriksa target existing: %v", err)
	}

	existingTargetMap := make(map[uint]*model.MarketingTargetBulanan)
	for _, et := range existingTargets {
		existingTargetMap[et.ProductID] = &et
	}

	for _, target := range req.Target {
		branchTarget, err := u.targetRepo.GetTargetBulananWithTx(tx, req.Tahun, req.Bulan, *bm.KantorCabangID, target.ProductID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("target bulanan cabang tidak ditemukan untuk produk %d", target.ProductID)
		}

		currentSum, err := u.targetRepo.GetMarketingTargetSumWithTx(tx, req.Tahun, req.Bulan, *bm.KantorCabangID, target.ProductID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal mendapatkan total target marketing: %v", err)
		}

		if existing, ok := existingTargetMap[target.ProductID]; ok {
			currentSum -= existing.TargetAmount
		}

		if currentSum+target.Amount > branchTarget.TargetAmount {
			tx.Rollback()
			return fmt.Errorf("total target marketing melebihi target bulanan cabang untuk produk %d", target.ProductID)
		}

		if existing, ok := existingTargetMap[target.ProductID]; ok {

			existing.TargetAmount = target.Amount
			if err := u.targetRepo.UpdateMarketingTargetWithTx(tx, existing); err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal mengupdate target marketing: %v", err)
			}
		} else {

			marketingTarget := &model.MarketingTargetBulanan{
				Tahun:                 req.Tahun,
				Bulan:                 req.Bulan,
				TargetAmount:          target.Amount,
				MarketingID:           marketing.ID,
				ProductID:             target.ProductID,
				KantorCabangID:        *bm.KantorCabangID,
				TargetProdukBulananID: branchTarget.ID,
			}
			if err := u.targetRepo.CreateMarketingTargetWithTx(tx, marketingTarget); err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal membuat target marketing: %v", err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menyimpan target: %v", err)
	}

	return nil
}

func (u *marketingTargetUsecase) GetMarketingTargets(req *dto.MonitoringRequest) ([]dto.MarketingTargetDetail, error) {

	result, err := u.targetRepo.GetMarketingTargets(req)
	return result, err
}
