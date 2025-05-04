package dto

type PredictionRequest struct {
	CIF                string   `json:"cif" validate:"required"`
	Nama               string   `json:"name" validate:"required"`
	NamaPerusahaan     string   `json:"company_name" validate:"required"`
	NomorRekening      string   `json:"nomor_rekening" validate:"required"`
	NomorHp            string   `json:"nomor_hp" validate:"required"`
	Alamat             string   `json:"address" validate:"required"`
	Pekerjaan          string   `json:"occupation" validate:"required"`
	Email              string   `json:"email" validate:"required,email"`
	Umur               int      `json:"umur" validate:"required"`
	Penghasilan        int64    `json:"income" validate:"required"`
	Payroll            bool     `json:"payroll"`
	Gender             string   `json:"gender" validate:"required"`
	StatusPerkawinan   bool     `json:"marital_status"`
	Segmen             string   `json:"category_segmen" validate:"required"`
	ProdukEksisting    []string `json:"existing_product" validate:"required"`
	AktivitasTransaksi string   `json:"transaction_activity" validate:"required,oneof=Active Inactive"`
}

type PredictionResult struct {
	Prediction string  `json:"prediksi"`
	Score      float64 `json:"skor"`
	Remarks    string  `json:"remarks"`
}
