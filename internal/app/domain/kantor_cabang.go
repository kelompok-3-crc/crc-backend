package dto

type CreateKantorCabangRequest struct {
	Nama string `json:"nama" validate:"required"`
}

type KantorCabangResponse struct {
	ID   uint   `json:"id"`
	Nama string `json:"nama"`
}
