package controller

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID    int    `json:"id"`
	CreatedAt int64  `json:"created_at"`
	ExpiresAt int64  `json:"expires_at"`
	Comment   string `json:"comment"`
	jwt.StandardClaims
}

func (s *Server) GenerateToken(user_id int, comment string) (string, error) {
	createdAt := time.Now().Unix() // data + sec = sec
	expiresAt := createdAt + 3600  // 1 hour from now

	claims := Claims{
		UserID:    user_id,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
		Comment:   comment,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}
