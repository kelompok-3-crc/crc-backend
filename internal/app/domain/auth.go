package dto

type LoginRequest struct {
	NIP      string `json:"nip" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateRequest struct {
	Nama     string `json:"nama" validate:"required,min=2"`
	Role     string `json:"role" validate:"required,oneof=bm user marketing"`
	Password string `json:"password" validate:"required,min=6"`
}
