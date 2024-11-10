package usecase

import (
	"crypto-project/internal/entity"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
)

func (u *Usecase) GetCurrencies(user_id int) (*entity.Response, error) {
	response, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,litecoin,tether&vs_currencies=usd")
	//формат: {"bitcoin":{"usd":68930},"ethereum":{"usd":2442.44},"litecoin":{"usd":66.82},"tether":{"usd":1.001}}
	if err != nil {
		u.logger.Error("http get error", slog.Any("error", err))
		return nil, errors.New("internal server error")
	}
	data, err := ioutil.ReadAll(response.Body)
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

func (u *Usecase) Convert(conversion *entity.Conversion) (float64, error) {
	response, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", conversion.To, conversion.From)) //conversion.From = usd
	//формат: {"tether":{"usd":1.001}}
	if err != nil {
		u.logger.Error("http get error", slog.Any("error", err))
		return 0, errors.New("internal server error")
	}
	data, err := ioutil.ReadAll(response.Body)
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
