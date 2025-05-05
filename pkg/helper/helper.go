package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func MapUnmarshalErrors(err error) map[string]string {
	errors := make(map[string]string)

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

	if invalidErr := strings.Split(err.Error(), "json: unknown field "); len(invalidErr) > 1 {
		fieldName := strings.Trim(invalidErr[1], "\"")
		errors[fieldName] = "Field tidak dikenali"
		return errors
	}

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

		const (
			DBR              = 0.51123123123
			MIN_TENOR_MONTHS = 12
			MAX_TENOR_MONTHS = 72
			MIN_MULTIPLIER   = 12
			MAX_MULTIPLIER   = 15 * 12
		)

		minPlafon := DBR * penghasilan * MIN_MULTIPLIER
		maxPlafon := DBR * penghasilan * MAX_MULTIPLIER

		if umur > 50 {
			reducedMultiplier := MAX_MULTIPLIER * (1.0 - float64(umur-50)/30.0)
			maxPlafon = DBR * penghasilan * max(MIN_MULTIPLIER, reducedMultiplier)
		}

		if payroll {
			maxPlafon *= 1.1
		}

		maxPlafon = min(maxPlafon, 1500000000)
		return Plafond{
			MinPlafon: uint64(minPlafon),
			MaxPlafon: uint64(maxPlafon),
			MinTenor:  MIN_TENOR_MONTHS,
			MaxTenor:  MAX_TENOR_MONTHS,
		}
	}

	if produk == "pensiun" {

		const DBR_PENSIUN = 0.9
		const PRICE_AKHIR_MIN = 11.5
		const PRICE_AKHIR_MAX = 15.0
		const RETIREMENT_AGE = 58
		const MAX_YEARS_BEFORE = 10

		tenorMaksAge := (75 - int(umur)) * 12
		tenorMaks15Years := 15 * 12
		tenorMaks := min(tenorMaksAge, tenorMaks15Years)

		angsuranMaks := DBR_PENSIUN * float64(penghasilan)

		priceAkhir := (PRICE_AKHIR_MIN + PRICE_AKHIR_MAX) / 2

		plafondMaks := angsuranMaks * priceAkhir

		yearsToRetirement := RETIREMENT_AGE - int(umur)

		eligibleForPrePension := (yearsToRetirement <= MAX_YEARS_BEFORE && yearsToRetirement > 0)
		if !eligibleForPrePension {
			plafondMaks = 0
			tenorMaks = 0
		}

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
		const MAX_AGE_AT_MATURITY = 70
		const MAX_PLAFOND = 5000000000

		angsuranMaks := float64(penghasilan) * DSR

		ageInMonths := int(umur) * 12
		maxTenorByAge := (MAX_AGE_AT_MATURITY * 12) - ageInMonths
		maxTenorActual := min(MAX_TENOR, maxTenorByAge)

		if maxTenorActual <= 0 {
			return Plafond{
				MinPlafon: 0,
				MaxPlafon: 0,
				MinTenor:  0,
				MaxTenor:  0,
			}
		}

		tenorBulanMax := float64(maxTenorActual)
		plafonMax := (angsuranMaks * tenorBulanMax) / (1 + (MARGIN * tenorBulanMax / 12))

		tenorBulanMin := float64(MIN_TENOR)
		plafonMin := (angsuranMaks * tenorBulanMin) / (1 + (MARGIN * tenorBulanMin / 12))

		if payroll {
			plafonMax *= 1.1
		}

		plafonMax = min(plafonMax, float64(MAX_PLAFOND))

		plafonMin = max(plafonMin, 100000000)

		return Plafond{
			MinPlafon: uint64(plafonMin),
			MaxPlafon: uint64(plafonMax),
			MinTenor:  MIN_TENOR,
			MaxTenor:  maxTenorActual,
		}
	}

	if produk == "oto" {

		const DSR = 0.5
		const MARGIN = 0.075
		const MIN_TENOR = 12
		const MAX_TENOR = 60
		const MAX_AGE_AT_MATURITY = 65
		const MAX_PLAFOND = 1000000000

		angsuranMaks := float64(penghasilan) * DSR

		ageInMonths := int(umur) * 12
		maxTenorByAge := (MAX_AGE_AT_MATURITY * 12) - ageInMonths
		maxTenorActual := min(MAX_TENOR, maxTenorByAge)

		if maxTenorActual <= 0 {
			return Plafond{
				MinPlafon: 0,
				MaxPlafon: 0,
				MinTenor:  0,
				MaxTenor:  0,
			}
		}

		tenorBulanMax := float64(maxTenorActual)
		pokokPembiayaanMax := (angsuranMaks * tenorBulanMax) / (1 + (MARGIN * tenorBulanMax / 12))

		tenorBulanMin := float64(MIN_TENOR)
		pokokPembiayaanMin := (angsuranMaks * tenorBulanMin) / (1 + (MARGIN * tenorBulanMin / 12))

		if payroll {
			pokokPembiayaanMax *= 1.05
		}

		pokokPembiayaanMax = min(pokokPembiayaanMax, float64(MAX_PLAFOND))

		return Plafond{
			MinPlafon: uint64(pokokPembiayaanMin),
			MaxPlafon: uint64(pokokPembiayaanMax),
			MinTenor:  MIN_TENOR,
			MaxTenor:  maxTenorActual,
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

		plafonMax = min(plafonMax, 250000000)

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
