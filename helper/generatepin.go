package helper

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePin(digit int) string {
	rand.Seed(time.Now().UnixNano())

	pin := ""
	for i := 1; i <= digit; i++ {
		pin += string(rune(rand.Intn(10) + '0'))
	}
	return pin
}

func VerifyOtp(phone, otp, otpHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(otpHash), []byte(otp))
	if err == nil {
		return true
	}
	return false
}

func IsNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
