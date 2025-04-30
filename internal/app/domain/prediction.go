package dto

type PredictionRequest struct {
	CIF                string   `json:"cif" validate:"required"`
	Nama               string   `json:"name" validate:"required"`
	NamaPerusahaan     string   `json:"company_name" validate:"required"`
	Umur               int      `json:"umur" validate:"required"`
	Penghasilan        int64    `json:"income" validate:"required"`
	Payroll            bool     `json:"payroll" validate:"required"`
	Gender             string   `json:"gender" validate:"required"`
	StatusPerkawinan   bool     `json:"marital_status"`
	Segmen             string   `json:"category_segmen" validate:"required"`
	ProdukEksisting    []string `json:"existing_product" validate:"required"`
	AktivitasTransaksi string   `json:"transaction_activity" validate:"required,oneof=active user inactive"`
}

type PredictionResult struct {
	Prediction string  `json:"prediksi"`
	Score      float64 `json:"skor"`
	Remarks    string  `json:"remarks"`
}
