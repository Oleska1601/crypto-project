package controller

import (
	"crypto-project/internal/entity"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log/slog"
	"net/http"
	"time"
)

// APIGetHandler страница ковертации
// @Summary get page
// @Description verify token and get currencies
// @Tags API
// @Accept json
// @Produce json
// @Param token header string true "jwt token for authentification"
// @Success 200 {string} string "Get currencies is successful"
// @Failure 401 {string} string "Authorization token not provided or Invalid token"
// @Failure 403 {string} string "Token has expired"
// @Failure 500 {string} string "Internal server error"
// @Router /api/get [get]
func (s *Server) APIGetHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("token")
	if tokenString == "" {
		http.Error(w, "Authorization token not provided", http.StatusUnauthorized)
		s.logger.Error("Authorization token not provided", slog.Int("status", http.StatusUnauthorized))
		return
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return s.secretKey, nil })
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		s.logger.Error("Invalid token", slog.Int("status", http.StatusUnauthorized))
		return
	}
	now := time.Now().Unix()
	if now > claims.ExpiresAt {
		http.Error(w, "Token has expired", http.StatusUnauthorized)
		s.logger.Error("Token has expired", slog.Int("status", http.StatusUnauthorized))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum,litecoin,tether&vs_currencies=usd")
	//формат: {"bitcoin":{"usd":68930},"ethereum":{"usd":2442.44},"litecoin":{"usd":66.82},"tether":{"usd":1.001}}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		s.logger.Error("get currencies error", slog.Any("error", err))
		return
	}
	result, err := s.u.GetCurrencies(claims.UserID, response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	s.logger.Info("Get currencies is successful", slog.Int("status", http.StatusOK))
	json.NewEncoder(w).Encode(result)
}

// APIConvertHandler страница конвертации
// @Summary convert page
// @Description verify token and convert amount
// @Tags API
// @Accept json
// @Produce json
// @Param token header string true "jwt token for authentification"
// @Param conversion body entity.Conversion true "amount, from, to"
// @Success 200 {string} string "convert is successful"
// @Failure 401 {string} string "Authorization token not provided or Invalid token"
// @Failure 403 {string} string "Token has expired"
// @Failure 500 {string} string "Internal server error"
// @Router /api/convert [post]
func (s *Server) APIConvertHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("token")
	if tokenString == "" {
		http.Error(w, "Authorization token not provided", http.StatusUnauthorized)
		s.logger.Error("Authorization token not provided", slog.Int("status", http.StatusUnauthorized))
		return
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return s.secretKey, nil })
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		s.logger.Error("Invalid token", slog.Int("status", http.StatusUnauthorized))
		return
	}
	now := time.Now().Unix()
	if now > claims.ExpiresAt {
		http.Error(w, "Token has expired", http.StatusForbidden)
		s.logger.Error("Token has expired", slog.Int("status", http.StatusForbidden))
		return
	}
	//return userID
	//json.NewEncoder(w).Encode(claims.UserID)

	var conversion entity.Conversion
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&conversion); err != nil {
		http.Error(w, "conversion is impossible", http.StatusBadRequest) //сервер не может обрабатывать запросы
		s.logger.Error("conversion error", slog.String("msg", "Bad Request"),
			slog.Int("status", http.StatusBadRequest), slog.Any("error", err))
		return
	}
	response, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", conversion.To, conversion.From)) //conversion.From = usd
	//формат: {"tether":{"usd":1.001}}
	result, err := s.u.Convert(&conversion, response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	s.logger.Info("convert is successful", slog.Int("status", http.StatusOK))
	json.NewEncoder(w).Encode(result)

}

// APIHistoryHandler страница истории
// @Summary history page
// @Description verify token and get history
// @Tags API
// @Accept json
// @Produce json
// @Param token header string true "jwt token for authentification"
// @Success 200 {string} string "Get history is successful"
// @Failure 401 {string} string "Authorization token not provided or Invalid token"
// @Failure 403 {string} string "Token has expired"
// @Failure 500 {string} string "Internal server error"
// @Router /api/history [get]
func (s *Server) APIHistoryHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("token")
	if tokenString == "" {
		http.Error(w, "Authorization token not provided", http.StatusUnauthorized)
		s.logger.Error("Authorization token not provided", slog.Int("status", http.StatusUnauthorized))
		return
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return s.secretKey, nil })
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		s.logger.Error("Invalid token", slog.Int("status", http.StatusUnauthorized))
		return
	}
	now := time.Now().Unix()
	if now > claims.ExpiresAt {
		http.Error(w, "Token has expired", http.StatusForbidden)
		s.logger.Error("Token has expired", slog.Int("status", http.StatusForbidden))
		return
	}

	result, err := s.u.GetHistory(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	s.logger.Info("Get history is successful", slog.Int("status", http.StatusOK))
	json.NewEncoder(w).Encode(result)
}
