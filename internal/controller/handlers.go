package controller

import (
	"crypto-project/internal/entity"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
)

// HomeHandler домашняя страница
// @Summary Home page
// @Description home page
// @Success 200
// @Failure 500 {string} string "Internal Server Error"
// @Router /home [get]
func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from homepage"))
}

// InfoHandler страница информации
// @Summary Info page
// @Description information about register and login
// @Success 200
// @Failure 500 {string} string "Internal Server Error"
// @Router /info [get]
func (s *Server) InfoHandler(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./web/info.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		s.logger.Error("Internal Server Error", slog.String("msg", "cannot parse files"),
			slog.Int("status", http.StatusInternalServerError), slog.Any("error", err))
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		s.logger.Error("Internal Server Error", slog.String("msg", "cannot execute"),
			slog.Int("status", http.StatusInternalServerError), slog.Any("error", err))
		return
	}
}

// LoginHandler страница аутентификации
// @Summary login page
// @Description login by json user with login, password
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.Credentials true "Login and password"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "login is impossible"
// @Failure 500 {string} string "login error"
// @Router /login [post]
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	//получаем только: login, password
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "login is impossible", http.StatusBadRequest)
		s.logger.Error("login error", slog.String("msg", "Bad Request"),
			slog.Int("status", http.StatusBadRequest), slog.Any("error", err))
		return
	}
	loginUser, err := s.u.LoginUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := s.GenerateToken(loginUser.ID, "simple token")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		s.logger.Error("GenerateToken error", slog.String("msg", "Internal Server Error"),
			slog.Int("status", http.StatusInternalServerError), slog.Any("error", err))
		return
	}
	w.Header().Set("token", token)
	w.WriteHeader(http.StatusOK)
	s.logger.Info("login successful", slog.String("user", loginUser.Login), slog.Int("status", http.StatusOK))
	//возращаем secret
	json.NewEncoder(w).Encode(loginUser.Secret)

}

// RegisterHandler страница регистрации
// @Summary registration page
// @Description register by json user with login, password
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.Credentials true "login and password"
// @Success 201 {string} string "registration successful"
// @Failure 400 {string} string "registration is impossible"
// @Failure 500 {string} string "registration error"
// @Router /register [post]
func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	//получаем только: login, password
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "registration is impossible", http.StatusBadRequest) //сервер не может обрабатывать запросы
		s.logger.Error("registration error", slog.String("msg", "Bad Request"),
			slog.Int("status", http.StatusBadRequest), slog.Any("error", err))
		return
	}
	if err := s.u.RegisterUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	s.logger.Info("registration successful", slog.String("user", user.Login), slog.Int("status", http.StatusCreated))
	//возращаем secret
	json.NewEncoder(w).Encode(user.Secret)
}
