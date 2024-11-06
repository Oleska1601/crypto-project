package usecase

import (
	"crypto-project/internal/entity"
	"encoding/base64"
	"errors"
	"log/slog"
)

func (u *Usecase) RegisterUser(user *entity.User) error {
	getUser, err := u.DB.GetTableUsers(user)
	if err != nil { //внутр ошибка
		u.logger.Error("GetTableUsers error", slog.Any("error", err))
		return errors.New("get user failed")
	}
	if UserExists(getUser) { //пользователь уже существует
		u.logger.Error("User already exists", slog.String("user", getUser.Login))
		return errors.New("incorrect login or password")
	}
	if err = u.InsertUser(user); err != nil {
		return errors.New("registration is failed")
	}
	return nil
}

func (u *Usecase) InsertUser(user *entity.User) error {
	user.Salt = base64.StdEncoding.EncodeToString(GenerateSalt())
	user.PasswordHash = GeneratePasswordHash(user.Password, user.Salt)
	user.Secret = GenerateSecret()
	err := u.DB.InsertTableUsers(user)
	if err != nil {
		u.logger.Error("InsertTableUsers error", slog.Any("error", err))
		return err
	}
	return nil
}

func (u *Usecase) LoginUser(user *entity.User) (*entity.User, error) {
	getUser, err := u.DB.GetTableUsers(user)
	if err != nil { //внутр ошибка
		u.logger.Error("GetTableUsers error", slog.Any("error", err))
		return nil, errors.New("get user failed")
	}
	if !UserExists(getUser) || !VerificationPassword(user.Password, getUser) {
		u.logger.Error("User not exists or password is incorrect", slog.String("user", user.Login))
		return nil, errors.New("incorrect login or password")
	}
	return &getUser, nil
}
