package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretkey1 = "jwt"
)

func New(id int, email string) (string, error) {

	const op = "utils.jwt.new"
	claim := jwt.RegisteredClaims{
		ID:        strconv.Itoa(id),
		Subject:   email,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	jwttoken, err := token.SignedString([]byte(secretkey1))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return jwttoken, nil

}
