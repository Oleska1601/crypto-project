package usecase

import (
	"crypto-project/internal/entity"
	"encoding/base64"
)

//отдельно вынесла функции, которые работают с бд - разумно ли их делать методами usecase???

func (u *Usecase) GetUser(user entity.User) (entity.User, error) {
	databaseUser, err := u.DB.GetTableUsers(user)
	return databaseUser, err
}

func (u *Usecase) RegisterUser(user *entity.User) error {
	user.Salt = base64.StdEncoding.EncodeToString(GenerateSalt())
	user.PasswordHash = GeneratePasswordHash(user.Password, user.Salt)
	user.Secret = GenerateSecret()
	err := u.DB.InsertTableUsers(user)
	return err
}
