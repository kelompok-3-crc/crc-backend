package dto

type UpdateCustomerStatusRequest struct {
	CIF       string `json:"cif" validate:"required,exists=customers.cif"`
	Status    string `json:"status" validate:"required,oneof=contacted rejected closed"`
	ProductID *uint  `json:"product_id" validate:"required_if=Status closed"`
	Amount    *int64 `json:"amount" validate:"required_if=Status closed,omitempty,min=1"`
	Notes     string `json:"notes" validate:"required_if=Status rejected"`
}
