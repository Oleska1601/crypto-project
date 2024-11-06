package usecase

import (
	"crypto-project/internal/entity"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log/slog"
)

func (u *Usecase) GetCurrencies(user_id int, body io.Reader) (*entity.Response, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		u.logger.Error("read body error", slog.Any("error", err))
		return nil, errors.New("internal server error")
	}
	var mapData map[string]map[string]float64
	if err = json.Unmarshal(data, &mapData); err != nil {
		u.logger.Error("Unmarshal data error", slog.Any("error", err))
		return nil, errors.New("internal server error")
	}
	result := &entity.Response{BTC: mapData["bitcoin"]["usd"], ETH: mapData["ethereum"]["usd"],
		LTC: mapData["litecoin"]["usd"], USDT: mapData["tether"]["usd"]}

	err = u.DB.InsertTableResponses(user_id, result)
	if err != nil {
		u.logger.Error("InsertTableResponses error", slog.Any("error", err))
		return nil, errors.New("internal server error")
	}
	return result, nil
}

func (u *Usecase) Convert(conversion *entity.Conversion, body io.Reader) (float64, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		u.logger.Error("read body error", slog.Any("error", err))
		return 0, errors.New("internal server error")
	}
	var mapData map[string]map[string]float64
	if err = json.Unmarshal(data, &mapData); err != nil {
		u.logger.Error("Unmarshal data error", slog.Any("error", err))
		return 0, errors.New("internal server error")
	}
	result := conversion.Amount / mapData[conversion.To][conversion.From]
	return result, nil
}

func (u *Usecase) GetHistory(userID int) ([]entity.Response, error) {
	result, err := u.DB.GetTableResponses(userID)
	if err != nil {
		u.logger.Error("GetTableResponses error", slog.Any("error", err))
		return nil, errors.New("internal server error")
	}
	return result, nil
}
