package helper

import (
	"math/rand"
	"time"
)

func GeneratePin(digit int) string {
	rand.Seed(time.Now().UnixNano())

	pin := ""
	for i := 1; i <= digit; i++ {
		pin += string(rune(rand.Intn(10) + '0'))
	}
	return pin
}
