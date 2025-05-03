package jwt

import (
	"ml-prediction/internal/app/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var accessSecret = []byte(os.Getenv("ACCESS_SECRET"))

type Claim struct {
	jwt.Claims
	NIP string `json:"nip"`
}

type JWTToken struct {
	Token    string
	Claim    Claim
	ExpireAt time.Time
	Scheme   string
}

func GenerateAccessToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"nip":  user.NIP,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}
