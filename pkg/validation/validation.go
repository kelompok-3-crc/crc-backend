package validation

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// MapValidationErrors maps validator errors into a field-to-message map.
func MapValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range verrs {
			// Use the JSON tag name if available
			field := e.Field()
			if jsonTag := e.StructField(); jsonTag != "" {
				field = strings.ToLower(jsonTag)
			}
			errors[strings.ToLower(field)] = getErrorMessage(e)
		}
	}
	return errors
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return fmt.Sprintf("Panjang minimal adalah %s karakter", e.Param())
	case "max":
		return fmt.Sprintf("Panjang maksimal adalah %s karakter", e.Param())
	case "len":
		return fmt.Sprintf("Panjang harus tepat %s karakter", e.Param())
	case "eq":
		return fmt.Sprintf("Harus sama dengan %s", e.Param())
	case "ne":
		return fmt.Sprintf("Tidak boleh sama dengan %s", e.Param())
	case "gte":
		return fmt.Sprintf("Harus lebih besar atau sama dengan %s", e.Param())
	case "gt":
		return fmt.Sprintf("Harus lebih besar dari %s", e.Param())
	case "lte":
		return fmt.Sprintf("Harus lebih kecil atau sama dengan %s", e.Param())
	case "lt":
		return fmt.Sprintf("Harus lebih kecil dari %s", e.Param())
	case "alphanum":
		return "Hanya boleh berisi huruf dan angka"
	case "numeric":
		return "Harus berupa angka"
	case "url":
		return "Harus berupa URL yang valid"
	case "uuid":
		return "Harus berupa UUID yang valid"
	case "oneof":
		return fmt.Sprintf("Harus salah satu dari [%s]", e.Param())
	default:
		return "Nilai tidak valid"
	}
}

func RegisterCustomValidation(v *validator.Validate) error {
	if err := v.RegisterValidation("fileformat", fileFormatValidator); err != nil {
		return fmt.Errorf("failed to register file format validation: %s", err)
	}

	if err := v.RegisterValidation("imageMaxSize", imageMaxSizeValidator); err != nil {
		return fmt.Errorf("failed to register image max size validation: %s", err)
	}

	if err := v.RegisterValidation("isBool", validateIsBool); err != nil {
		return fmt.Errorf("failed to register boolean validation: %s", err)
	}

	if err := v.RegisterValidation("noSpace", validateNoSpace); err != nil {
		return fmt.Errorf("failed to register username has space: %s", err)
	}

	return nil
}

func CustomError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "min", "max":
		return fmt.Sprintf("%s too short or long", e.Field())
	default:
		return e.Error()
	}
}

func ValidateFile(v *validator.Validate, fileHeader *multipart.FileHeader) error {

	if err := v.Var(fileHeader, "fileformat"); err != nil {
		return fmt.Errorf("file format must be JPG or JPEG")
	}
	if err := v.Var(fileHeader, "imageMaxSize"); err != nil {
		return fmt.Errorf("image size cannot exceed 2MB")
	}
	return nil
}

func fileFormatValidator(fl validator.FieldLevel) bool {
	file, ok := fl.Top().Interface().(*multipart.FileHeader)
	if !ok {
		return false
	}

	ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, ".")+1:])
	return ext == "jpg" || ext == "jpeg"
}

func imageMaxSizeValidator(fl validator.FieldLevel) bool {
	file, ok := fl.Top().Interface().(*multipart.FileHeader)
	if !ok {
		return false
	}
	maxSize := int64(2) // 2MB
	return file.Size <= maxSize
}

func validateIsBool(fl validator.FieldLevel) bool {
	return fl.Field().Kind() == reflect.Bool
}

func validateNoSpace(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return !strings.Contains(field, " ")
}

func UrlValidation(url string) error {
	pattern := `^(http(s)?:\/\/)[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}(\/?([-a-zA-Z0-9@:%_\+.~#?&//=]*\.(png|jpg|jpeg|gif)))?$`

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile pattern %v", err)
	}

	if !regex.MatchString(url) {
		return fmt.Errorf("url is not valid")
	}

	return nil
}

func UuidValidation(uuid string) error {
	pattern := `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile pattern %v", err)
	}

	if !regex.MatchString(uuid) {
		return fmt.Errorf("uuid is not valid")
	}

	return nil
}

func EmailValidation(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile pattern %v", err)
	}

	if !regex.MatchString(email) {
		return fmt.Errorf("email is not valid")
	}

	return nil
}

func PhoneValidation(phone string) error {
	if !strings.HasPrefix(phone, "+") {
		return fmt.Errorf("phone is not valid")
	}

	pattern := `^\+\d+$`
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(phone) {
		return fmt.Errorf("phone is not valid")
	}

	if len(phone) < 7 || len(phone) > 14 {
		return fmt.Errorf("phone is too short or long")
	}

	return nil
}

func ValidateImageFileType(fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:])
	if !(ext == "jpg" || ext == "jpeg") {
		return fmt.Errorf("file format must be JPG or JPEG")
	}
	return nil
}

func ValidateParams(r *http.Request, fields interface{}) error {

	val := reflect.ValueOf(fields)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i).Tag.Get("json")
		if _, exist := r.Form[field]; exist {
			if r.Form.Get(field) == "" {
				return fmt.Errorf("%s field should have value if present", field)
			}
		}
	}
	return nil
}
