package usecase

import (
	"crypto-project/internal/entity"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func UserExists(user entity.User) bool {
	return user != entity.User{}
}

func GenerateSalt() []byte {
	saltBytes := make([]byte, 16)
	rand.Read(saltBytes)
	return saltBytes
}

func GeneratePasswordHash(password string, salt string) string {
	passHash := sha256.Sum256(append([]byte(password), []byte(salt)...))
	return base64.StdEncoding.EncodeToString(passHash[:])
}

func GenerateSecret() string {
	secretBytes := make([]byte, 16)
	rand.Read(secretBytes)
	return base64.StdEncoding.EncodeToString(secretBytes)
}

func VerificationPassword(inputPass string, user entity.User) bool {
	inputPassHash := GeneratePasswordHash(inputPass, user.Salt)
	return user.PasswordHash == inputPassHash
}
