package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	dto "main/internal/dto"
	"main/package/util"

	"github.com/golang-jwt/jwt/v4"
)

var (
	JWT_SECRET         = []byte(util.Getenv("JWT_SECRET", "testsecret"))
	JWT_EXP            = time.Duration(168) * time.Hour
	JWT_SIGNING_METHOD = jwt.SigningMethodHS256
)

func getTokenString(authHeader string) (*string, error) {
	var token string
	if strings.Contains(authHeader, "Bearer") {
		token = strings.Replace(authHeader, "Bearer ", "", -1)
		return &token, nil
	}
	return nil, fmt.Errorf("authorization not found")
}

func CreateJWTClaims(uuid, email, phone, roles string, permissions []string, admin bool) dto.JWTClaims {
	return dto.JWTClaims{
		Uuid:        uuid,
		Email:       email,
		Phone:       phone,
		Roles:       roles,
		Permissions: permissions,
		Admin:       admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWT_EXP)),
		},
	}
}

func CreateJWTToken(claims dto.JWTClaims) (string, error) {
	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)
	return token.SignedString([]byte(JWT_SECRET))
}

func ParseJWTToken(authHeader string) (*dto.JWTClaims, error) {
	tokenString, err := getTokenString(authHeader)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method != JWT_SIGNING_METHOD {
			return nil, fmt.Errorf("invalid signing method")
		}
		return JWT_SECRET, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claimsStr, err := json.Marshal(claims)
		if err != nil {
			return nil, fmt.Errorf("error when marshalling token")
		}

		var customClaims dto.JWTClaims
		if err := json.Unmarshal(claimsStr, &customClaims); err != nil {
			return nil, fmt.Errorf("error when unmarshalling token")
		}

		return &customClaims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
