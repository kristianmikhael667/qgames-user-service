package helper

import (
	"fmt"
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
	fmt.Println("otp ", otp)
	fmt.Println("has ", otpHash)
	err := bcrypt.CompareHashAndPassword([]byte(otpHash), []byte(otp))
	if err == nil {
		return true
	}
	return false
}
