package dto

// TargetProdukTahunanRequest contains data for yearly target creation/update.
// @Description Yearly target request
type TargetProdukTahunanRequest struct {
	// The year for which targets are being set
	// @example 2025
	Tahun int `json:"tahun" validate:"required"`

	// List of product targets
	Targets []ProductTarget `json:"target" validate:"required,min=1,all_products,dive"`
}

// TargetProdukRequest contains target data for a specific product.
// @Description Target for a specific product
type ProductTarget struct {
	// Product ID
	// @example 1
	ProductID   uint   `json:"product_id" validate:"required"`
	ProductName string `json:"product_name" gorm:"column:product_name"`

	// Target amount for the year
	// @example 1000000000
	Amount int64 `json:"amount" validate:"required,min=0"`
}

type AssignMarketingTargetRequest struct {
	Tahun        int             `json:"tahun" validate:"required,min=2024"`
	Bulan        int             `json:"bulan" validate:"required,min=1,max=12"`
	MarketingNIP string          `json:"marketing_nip" validate:"required"`
	Target       []ProductTarget `json:"target" validate:"required,min=1,all_products,dive"`
}

type SingleMarketingTarget struct {
	Tahun        int    `json:"tahun" validate:"required,min=2024"`
	Bulan        int    `json:"bulan" validate:"required,min=1,max=12"`
	MarketingNIP string `json:"marketing_nip" validate:"required"`
	ProductID    uint   `json:"product_id" validate:"required"`
	Amount       int64  `json:"amount" validate:"required,min=1"`
}

type ProductTargetSummary struct {
	// Product ID
	// @example 1
	ProductID uint `json:"product_id"`
	// Product name
	// @example "Hasanah Card"
	ProductName string `json:"product_name"`
	// Product target amount
	// @example 5000000000
	Target float64 `json:"target"`
	// Achieved amount for the product
	// @example 3750000000
	Achieved float64 `json:"achieved"`
	// Achievement percentage for the product
	// @example 75.0
	Percentage float64 `json:"percentage"`
}

type TargetSummaryResponse struct {
	Type       string `json:"type"` // "branch" or "marketing"
	BranchName string `json:"branch_name" gorm:"column:branch_name"`
	Name       string `json:"name" gorm:"column:name"`
	NIP        string `json:"nip,omitempty"`
	// Total target amount across all products
	// @example 10000000000
	TotalTarget float64 `json:"total_target"`
	// Total achievement amount across all products
	// @example 7550000000
	Achieved float64 `json:"achieved"`
	// Overall branch target achievement percentage
	// @example 75.5
	Percentage float64 `json:"percentage"`
	// List of product-specific targets and achievements
	Products     []ProductTargetSummary `json:"products"`
	TargetMonth  int                    `json:"target_month"`
	TargetYear   int                    `json:"target_year"`
	TargetSetted bool                   `json:"target_setted"`
}

type BranchProductTarget struct {
	ProductID        uint    `json:"product_id"`
	ProductName      string  `json:"product_name"`
	TotalTarget      float64 `json:"total_target"`
	AssignedAmount   float64 `json:"assigned_amount"`
	UnassignedAmount float64 `json:"unassigned_amount"`
}

type BranchTargetResponse struct {
	BranchID   uint                  `json:"branch_id"`
	BranchName string                `json:"branch_name"`
	Month      int                   `json:"month"`
	Year       int                   `json:"year"`
	Products   []BranchProductTarget `json:"products"`
}
