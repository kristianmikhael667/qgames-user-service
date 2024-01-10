package helper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func SendOtp(phone string, otp string) (string, int) {
	url := os.Getenv("VENDOR_QONTAK")
	method := "POST"
	phoneNumber := phone
	phoneNumber = strings.TrimPrefix(phoneNumber, "0")
	phoneNumber = "62" + phoneNumber

	payload := map[string]interface{}{
		"to_number":              phoneNumber,
		"to_name":                "users_qgrowid",
		"message_template_id":    os.Getenv("MESSAGE_TEMPLATE"),
		"channel_integration_id": os.Getenv("CHANNEL_ID"),
		"language": map[string]interface{}{
			"code": "id",
		},
		"parameters": map[string]interface{}{
			"body": []map[string]interface{}{
				{
					"key":        "1",
					"value":      "send_otp",
					"value_text": otp,
				},
			},
		},
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Print("Error payloadBytes : ", err)
		return "Error payloadBytes", 500
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Print("Error HTTP Create request : ", err)
		return "Error HTTP Create request", 500
	}

	// Set request headers
	req.Header.Set("Authorization", os.Getenv("TOKEN_QONTAK"))
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Print("Error Send HTTP request : ", err, " sc ? ", res.StatusCode)
		return "Error Send HTTP request", res.StatusCode
	}
	defer res.Body.Close()

	// Process response
	log.Print("Response status:", res.Status)
	log.Print("Success Send OTP to number: "+phone, res.Status)
	return res.Status, res.StatusCode
}
