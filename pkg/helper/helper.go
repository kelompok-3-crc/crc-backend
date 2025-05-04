package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func MapUnmarshalErrors(err error) map[string]string {
	errors := make(map[string]string)

	// Check if it's a JSON syntax error
	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		errors["json"] = fmt.Sprintf("JSON syntax error at position %d: %v", syntaxErr.Offset, syntaxErr.Error())
		return errors
	}

	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		errors[typeErr.Field] = fmt.Sprintf(
			"Field harus bertipe %v, bukan %v",
			typeErr.Type.String(),
			typeErr.Value,
		)
		return errors
	}

	// Handle invalid field names
	if invalidErr := strings.Split(err.Error(), "json: unknown field "); len(invalidErr) > 1 {
		fieldName := strings.Trim(invalidErr[1], "\"")
		errors[fieldName] = "Field tidak dikenali"
		return errors
	}

	// Default error
	errors["json"] = err.Error()
	return errors
}
func GetProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return wd, nil
}

type Plafond struct {
	MinPlafon uint64
	MaxPlafon uint64
	MinTenor  int
	MaxTenor  int
}

func CalculatePlafond(produk string, umur, penghasilan int64, payroll bool) Plafond {
	if produk == "mitraguna" {
		penghasilan := float64(penghasilan)
		minPlafon := 0.7 * penghasilan * 12
		maxPlafon := 0.7 * penghasilan * 15 * 12
		minTenor := 12
		maxTenor := 15 * 12

		return Plafond{
			MinPlafon: uint64(minPlafon),
			MaxPlafon: uint64(maxPlafon),
			MinTenor:  minTenor,
			MaxTenor:  maxTenor,
		}
	}
	if produk == "pensiun" {

		const DBR_PENSIUN = 0.9
		const PRICE_AKHIR_MIN = 11.5
		const PRICE_AKHIR_MAX = 15.0

		tenorMaksAge := (75 - int(umur)) * 12
		tenorMaks15Years := 15 * 12
		tenorMaks := min(tenorMaksAge, tenorMaks15Years)

		angsuranMaks := DBR_PENSIUN * float64(penghasilan)

		priceAkhir := (PRICE_AKHIR_MIN + PRICE_AKHIR_MAX) / 2

		plafondMaks := angsuranMaks * priceAkhir

		return Plafond{
			MinPlafon: 0,
			MaxPlafon: uint64(plafondMaks),
			MinTenor:  0,
			MaxTenor:  tenorMaks,
		}
	}
	if produk == "prapensiun" {
		const DBR_PRA_PENSIUN = 0.4
		const PRICE_AKHIR_MIN = 12.0
		const PRICE_AKHIR_MAX = 13.0
		const RETIREMENT_AGE = 58
		const MAX_YEARS_BEFORE = 10
		const MAX_TENOR_YEARS = 15

		yearsToRetirement := RETIREMENT_AGE - int(umur)

		eligibleForPrePension := (yearsToRetirement <= MAX_YEARS_BEFORE && yearsToRetirement > 0)

		maxTenorByAge := (75 - int(umur)) * 12
		maxTenorByYears := MAX_TENOR_YEARS * 12
		tenorMaks := min(maxTenorByAge, maxTenorByYears)

		angsuranMaks := DBR_PRA_PENSIUN * float64(penghasilan)

		priceAkhir := (PRICE_AKHIR_MIN + PRICE_AKHIR_MAX) / 2

		plafondMaks := angsuranMaks * priceAkhir

		if !eligibleForPrePension {
			plafondMaks = 0
			tenorMaks = 0
		}

		tenorBeforePension := min(yearsToRetirement*12, tenorMaks)
		tenorAfterPension := max(0, tenorMaks-tenorBeforePension)

		return Plafond{
			MinPlafon: 0,
			MaxPlafon: uint64(plafondMaks),
			MinTenor:  tenorAfterPension + tenorBeforePension,
			MaxTenor:  tenorMaks,
		}
	}

	if produk == "griya" {
		const DSR = 0.4
		const MARGIN = 0.1
		const MIN_TENOR = 12
		const MAX_TENOR = 300

		angsuranMaks := float64(penghasilan) * DSR

		tenorBulanMax := float64(MAX_TENOR)
		plafonMax := (angsuranMaks * tenorBulanMax) / (1 + (MARGIN * tenorBulanMax / 12))

		tenorBulanMin := float64(MIN_TENOR)
		plafonMin := (angsuranMaks * tenorBulanMin) / (1 + (MARGIN * tenorBulanMin / 12))

		return Plafond{
			MinPlafon: uint64(plafonMin),
			MaxPlafon: uint64(plafonMax),
			MinTenor:  MIN_TENOR,
			MaxTenor:  MAX_TENOR,
		}
	}

	if produk == "oto" {
		const DSR = 0.5
		const MARGIN = 0.075
		const MIN_TENOR = 12
		const MAX_TENOR = 60

		angsuranMaks := float64(penghasilan) * DSR

		tenorBulanMax := float64(MAX_TENOR)
		pokokPembiayaanMax := (angsuranMaks * tenorBulanMax) / (1 + (MARGIN * tenorBulanMax / 12))

		tenorBulanMin := float64(MIN_TENOR)
		pokokPembiayaanMin := (angsuranMaks * tenorBulanMin) / (1 + (MARGIN * tenorBulanMin / 12))

		return Plafond{
			MinPlafon: uint64(pokokPembiayaanMin),
			MaxPlafon: uint64(pokokPembiayaanMax),
			MinTenor:  MIN_TENOR,
			MaxTenor:  MAX_TENOR,
		}
	}
	if produk == "hasanahcard" {
		var faktorLimit float64

		if umur < 35 {
			faktorLimit = 2.5
		} else if umur >= 35 && umur < 45 {

			faktorLimit = 2.0
		} else if umur >= 45 && umur < 60 {

			faktorLimit = 1.5
		} else {
			faktorLimit = 1.0
		}

		plafonMax := float64(penghasilan) * faktorLimit

		plafonMin := float64(penghasilan) * 1.0

		return Plafond{
			MinPlafon: uint64(plafonMin),
			MaxPlafon: uint64(plafonMax),
			MinTenor:  0,
			MaxTenor:  0,
		}
	}
	return Plafond{
		MinPlafon: 0,
		MaxPlafon: 0,
		MinTenor:  0,
		MaxTenor:  0,
	}
}
