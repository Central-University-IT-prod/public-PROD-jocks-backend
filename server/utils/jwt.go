package utils

import (
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(payload any, secret string) (string, error) {
	claims, err := ToMap(payload)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ReadJWT[T any](item *T, tokenString, secret string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid payload")
	}

	data, err := json.Marshal(claims)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, item)
}
