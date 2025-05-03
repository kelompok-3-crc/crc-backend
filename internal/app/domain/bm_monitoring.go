package dto

type MarketingMonitoringResponse struct {
	MarketingNIP    string         `json:"marketing_nip" gorm:"column:marketing_nip"`
	MarketingName   string         `json:"marketing_name" gorm:"column:marketing_name"`
	MonthlyAchieved float64        `json:"monthly_achieved" gorm:"column:monthly_achieved"`
	MonthlyTarget   float64        `json:"monthly_target" gorm:"column:monthly_target"`
	Labels          []string       `json:"labels" gorm:"-"`         // Ignore in GORM as it's for ChartJS
	Datasets        []ChartDataset `json:"datasets" gorm:"-"`       // Ignore in GORM as it's for ChartJS
	TargetDetails   []ProductChart `json:"target_details" gorm:"-"` // Ignore in GORM as it's for ChartJS
}

type ChartDataset struct {
	Data            []float64 `json:"data"`
	BackgroundColor []string  `json:"backgroundColor"`
	BorderWidth     int       `json:"borderWidth"`
}
type ProductChart struct {
	ProductID       uint           `json:"product_id" gorm:"column:product_id"`
	ProductName     string         `json:"product_name" gorm:"column:product_name"`
	MonthlyAchieved float64        `json:"monthly_achieved" gorm:"column:monthly_achieved"`
	MonthlyTarget   float64        `json:"monthly_target" gorm:"column:monthly_target"`
	Labels          []string       `json:"labels" gorm:"-"`   // Ignore in GORM as it's for ChartJS
	Datasets        []ChartDataset `json:"datasets" gorm:"-"` // Ignore in GORM as it's for ChartJS
}

type MonitoringRequest struct {
	Month     int    `json:"month" validate:"required,min=1,max=12"`
	Year      int    `json:"year" validate:"required,min=2024"`
	HasTarget *bool  `json:"has_target,omitempty"`
	Search    string `json:"search,omitempty"`
	NIP       string `json:"nip,omitempty"`
}

type MarketingTargetDetail struct {
	MarketingNIP  string          `json:"marketing_nip" gorm:"column:marketing_nip"`
	MarketingName string          `json:"marketing_name" gorm:"column:marketing_name"`
	HasTarget     bool            `json:"has_target" gorm:"column:has_target"`
	TotalTarget   float64         `json:"total_target" gorm:"column:total_target"`
	TargetDetails []ProductTarget `json:"target_details"`
}
type ProductPerformanceRequest struct {
	ProductID uint   `json:"product_id" form:"product_id" validate:"required"`
	StartDate string `json:"start_date" form:"start_date" validate:"required"`
	EndDate   string `json:"end_date" form:"end_date" validate:"required"`
	GroupBy   string `json:"group_by" form:"group_by" validate:"required,oneof=week month year"`
}
type ProductPerformanceResponse struct {
	Produk   string             `json:"produk"`
	ProdukID string             `json:"produk_id"`
	Labels   []string           `json:"labels"`
	Datasets []MarketingDataset `json:"datasets"`
}

type MarketingDataset struct {
	MarketingNIP    string    `json:"marketing_nip"`
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	Targets         []float64 `json:"targets"` // Per-period targets
	BorderColor     string    `json:"borderColor"`
	BackgroundColor string    `json:"backgroundColor"`
	Fill            bool      `json:"fill"`
	Tension         float64   `json:"tension"`
	Target          float64   `json:"target"` // Annual target (total)
}
