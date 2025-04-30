package utils

import (
	"math/rand"
	"time"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomNIP(length int) string {
	digits := "0123456789"
	nip := make([]byte, length)
	for i := 0; i < length; i++ {
		nip[i] = digits[seededRand.Intn(len(digits))]
	}
	return "M" + string(nip)
}
