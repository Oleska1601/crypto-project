package controller

import (
	"crypto-project/internal/entity"
	"encoding/json"
	"log/slog"
	"net/http"
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
	userID := r.Context().Value("userID").(int)
	result, err := s.u.GetCurrencies(userID)
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
	var conversion entity.Conversion
	if err := json.NewDecoder(r.Body).Decode(&conversion); err != nil {
		http.Error(w, "conversion is impossible", http.StatusBadRequest) //сервер не может обрабатывать запросы
		s.logger.Error("conversion error", slog.String("msg", "Bad Request"),
			slog.Int("status", http.StatusBadRequest), slog.Any("error", err))
		return
	}
	result, err := s.u.Convert(&conversion)
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
	userID := r.Context().Value("userID").(int)
	result, err := s.u.GetHistory(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	s.logger.Info("Get history is successful", slog.Int("status", http.StatusOK))
	json.NewEncoder(w).Encode(result)
}
