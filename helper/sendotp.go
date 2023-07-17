package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func SendOtp(phone string, otp string) (string, bool) {
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
		fmt.Println("Error:", err)
		return err.Error(), false
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error(), false
	}

	// Set request headers
	req.Header.Set("Authorization", os.Getenv("TOKEN_QONTAK"))
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error(), false
	}
	defer res.Body.Close()

	// Process response
	fmt.Println("Response status:", res.Status)
	Logger("info", "Success Send OTP to number: "+phone, "Rc: "+res.Status)
	return res.Status, true
}
