// Package usecase contains business logic for the application.
// @title BSI Prediction API
// @version 1.0
// @description API for BSI Prediction and Target Management
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@bsi.co.id
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
package usecase

import (
	"context"
	"fmt"
	dto "ml-prediction/internal/app/domain"
	"ml-prediction/internal/app/model"
	"ml-prediction/internal/app/repository"

	"gorm.io/gorm"
)

// TargetUsecase handles business logic for target management
// @Summary Interface for target-related operations
type TargetUsecase interface {
	// CreateOrUpdateTargetTahunan creates or updates yearly targets and distributes them to monthly targets
	// @Summary Creates or updates yearly targets
	// @Description Sets yearly targets for products and automatically distributes them as monthly targets
	// @Tags targets
	// @Accept json
	// @Produce json
	// @Param userNIP path string true "User NIP"
	// @Param request body dto.TargetProdukTahunanRequest true "Target yearly request"
	// @Success 200 {object} string "Target berhasil disimpan"
	// @Failure 400 {object} string "Validation error message"
	// @Failure 401 {object} string "Unauthorized"
	// @Failure 500 {object} string "Internal server error"
	// @Router /targets/yearly [post]
	CreateOrUpdateTargetTahunan(userNIP string, req *dto.TargetProdukTahunanRequest) error

	// GetTargetSummary retrieves target summary for a specific month and year
	// @Summary Gets target summary
	// @Description Returns target summary including achievements for the specified month and year
	// @Tags targets
	// @Accept json
	// @Produce json
	// @Param userID path int true "User ID"
	// @Param userNIP path string true "User NIP"
	// @Param month query int true "Month number (1-12)"
	// @Param year query int true "Year (e.g. 2025)"
	// @Success 200 {object} dto.TargetSummaryResponse "Target summary"
	// @Failure 401 {object} string "Unauthorized"
	// @Failure 404 {object} string "User not found"
	// @Failure 500 {object} string "Internal server error"
	// @Router /targets/summary [get]
	GetTargetSummary(ctx context.Context, userID uint, userNIP string, month, year int) (*dto.TargetSummaryResponse, error)

	// GetBranchMonthlyTargetWithUnassigned retrieves branch monthly targets with unassigned amounts
	GetBranchMonthlyTargetWithUnassigned(ctx context.Context, userNIP string, month, year int) (*dto.BranchTargetResponse, error)
}

type targetUsecase struct {
	userRepo   repository.UserRepository
	targetRepo repository.TargetRepository
	db         *gorm.DB
}

func NewTargetUsecase(userRepo repository.UserRepository, targetRepo repository.TargetRepository,
	db *gorm.DB) TargetUsecase {
	return &targetUsecase{
		userRepo:   userRepo,
		targetRepo: targetRepo,
		db:         db,
	}
}

func (u *targetUsecase) CreateOrUpdateTargetTahunan(userNIP string, req *dto.TargetProdukTahunanRequest) error {
	tx := u.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := u.userRepo.FindByNIPWithTx(tx, userNIP)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("user tidak ditemukan")
	}
	if user.Role != "bm" {
		tx.Rollback()
		return fmt.Errorf("hanya Regional Manager yang dapat mengatur target tahunan")
	}

	existingTargets, err := u.targetRepo.GetExistingTargetTahunanWithTx(tx, req.Tahun, *user.KantorCabangID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal memeriksa target existing: %v", err)
	}

	existingTargetMap := make(map[uint]*model.TargetProdukTahunan)
	for _, et := range existingTargets {
		existingTargetMap[et.ProductID] = &et
	}

	for _, target := range req.Targets {
		var tahunanID uint
		if existing, ok := existingTargetMap[target.ProductID]; ok {
			existing.TargetAmount = target.Amount
			if err := u.targetRepo.UpdateTargetTahunanWithTx(tx, existing); err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal mengupdate target tahunan: %v", err)
			}
			tahunanID = existing.ID
		} else {
			newTarget := &model.TargetProdukTahunan{
				Tahun:          req.Tahun,
				TargetAmount:   target.Amount,
				KantorCabangID: *user.KantorCabangID,
				ProductID:      target.ProductID,
			}
			if err := u.targetRepo.CreateTargetTahunanWithTx(tx, newTarget); err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal membuat target tahunan: %v", err)
			}
			tahunanID = newTarget.ID
		}

		// Calculate monthly target (divide yearly target by 12)
		monthlyAmount := target.Amount / 12

		// Update or create monthly targets for all months
		for month := 1; month <= 12; month++ {
			existingMonthly, err := u.targetRepo.GetTargetBulananWithTx(tx, req.Tahun, month, *user.KantorCabangID, target.ProductID)
			if err != nil && err != gorm.ErrRecordNotFound {
				tx.Rollback()
				return fmt.Errorf("gagal memeriksa target bulanan: %v", err)
			}

			if existingMonthly != nil {
				existingMonthly.TargetAmount = monthlyAmount
				if err := u.targetRepo.UpdateTargetBulananWithTx(tx, existingMonthly); err != nil {
					tx.Rollback()
					return fmt.Errorf("gagal mengupdate target bulanan: %v", err)
				}
			} else {
				newMonthly := &model.TargetProdukBulanan{
					Tahun:           req.Tahun,
					Bulan:           month,
					TargetAmount:    monthlyAmount,
					KantorCabangID:  *user.KantorCabangID,
					ProductID:       target.ProductID,
					TargetTahunanID: tahunanID,
				}
				if err := u.targetRepo.CreateTargetBulananWithTx(tx, newMonthly); err != nil {
					tx.Rollback()
					return fmt.Errorf("gagal membuat target bulanan: %v", err)
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal menyimpan target: %v", err)
	}

	return nil
}

func (u *targetUsecase) GetTargetSummary(ctx context.Context, userID uint, userNIP string, month, year int) (*dto.TargetSummaryResponse, error) {
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user model.User
	if err := tx.Where("id = ? AND nip = ? AND deleted_at IS NULL", userID, userNIP).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	summary, err := u.targetRepo.GetTargetSummary(tx, userID, user.Role, month, year)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return summary, nil
}

func (u *targetUsecase) GetBranchMonthlyTargetWithUnassigned(ctx context.Context, userNIP string, month, year int) (*dto.BranchTargetResponse, error) {
	user, err := u.userRepo.FindByNIP(userNIP)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	if user.Role != "bm" {
		return nil, fmt.Errorf("only branch managers can access branch targets")
	}

	// Get branch targets for the month
	branchTargets, err := u.targetRepo.GetBranchTargets(*user.KantorCabangID, month, year)
	if err != nil {
		return nil, fmt.Errorf("error getting branch targets: %v", err)
	}

	// Get assigned targets to marketing staff
	assignedTargets, err := u.targetRepo.GetAssignedTargets(*user.KantorCabangID, month, year)
	if err != nil {
		return nil, fmt.Errorf("error getting assigned targets: %v", err)
	}

	// Calculate unassigned amounts
	response := &dto.BranchTargetResponse{
		BranchID:   *user.KantorCabangID,
		BranchName: "", // Fetch branch name if needed
		Month:      month,
		Year:       year,
		Products:   []dto.BranchProductTarget{},
	}

	// Create a map of assigned amounts by product
	assignedMap := make(map[uint]float64)
	for _, assigned := range assignedTargets {
		assignedMap[assigned.ProductID] += assigned.Amount
	}

	// Calculate unassigned amounts for each product
	for _, target := range branchTargets {
		assigned := assignedMap[target.ProductID]
		unassigned := target.TargetAmount - assigned
		if unassigned < 0 {
			unassigned = 0 // Cannot have negative unassigned amount
		}

		response.Products = append(response.Products, dto.BranchProductTarget{
			ProductID:        target.ProductID,
			ProductName:      target.ProductName,
			TotalTarget:      target.TargetAmount,
			AssignedAmount:   assigned,
			UnassignedAmount: unassigned,
		})
	}

	return response, nil
}
