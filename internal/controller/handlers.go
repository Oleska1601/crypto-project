package controller

import (
	"crypto-project/internal/entity"
	"crypto-project/internal/usecase"
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
// @Param user body entity.User true "Login and password"
// @Success 200 {string} string ""Login successful"
// @Failure 400 {string} string "login is impossible"
// @Failure 500 {string} string "login error"
// @Router /login [post]
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	w.Header().Set("Content-Type", "application/json")
	//получаем только: login, password
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "login is impossible", http.StatusBadRequest)
		s.logger.Error("login error", slog.String("msg", "Bad Request"),
			slog.Int("status", http.StatusBadRequest), slog.Any("error", err))
		return
	}
	//getUser = пустой структуре, если user не существует (проверка на логин),
	//в противном случае получаем его login, password_hash, salt, secret
	getUser, err := s.u.GetUser(user)
	if err != nil { //какая-то внутренняя ошибка
		http.Error(w, "get user failed", http.StatusInternalServerError) //сервер не может обрабатывать запросы
		s.logger.Error("login error", slog.String("msg", "get user failed"),
			slog.Int("status", http.StatusInternalServerError), slog.Any("error", err))
		return
	}
	if !usecase.UserExists(getUser) {
		http.Error(w, "user does not exist", http.StatusInternalServerError) //сервер не может обрабатывать запросы
		s.logger.Error("login error", slog.String("msg", "user does not exist"),
			slog.Int("status", http.StatusInternalServerError), slog.Any("error", err))
		return
	}
	if !usecase.VerificationPassword(user.Password, getUser) {
		http.Error(w, "incorrect password", http.StatusInternalServerError) //сервер не может обрабатывать запросы
		s.logger.Error("login error", slog.String("msg", "incorrect password"),
			slog.Int("status", http.StatusInternalServerError), slog.Any("error", err))
		return
	}
	w.WriteHeader(http.StatusOK)
	s.logger.Info("login successful", slog.Int("status", http.StatusOK))
	//возращаем secret
	json.NewEncoder(w).Encode(getUser.Secret)

}

// RegisterHandler страница регистрации
// @Summary registration page
// @Description register by json user with login, password
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.User true "login and password"
// @Success 201 {string} string "registration successful"
// @Failure 400 {string} string "register is impossible"
// @Failure 500 {string} string "register error"
// @Router /register [post]
func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	w.Header().Set("Content-Type", "application/json")
	//получаем только: login, password
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "registration is impossible", http.StatusBadRequest) //сервер не может обрабатывать запросы
		s.logger.Error("registration error", slog.String("msg", "Bad Request"),
			slog.Int("status", http.StatusBadRequest), slog.Any("error", err))
		return
	}
	//getUser = пустой структуре, если user не существует (проверка на логин),
	//в противном случае получаем его login, password_hash, salt, secret
	getUser, err := s.u.GetUser(user)
	if err != nil {
		http.Error(w, "get user failed", http.StatusInternalServerError) //сервер не может обрабатывать запросы
		s.logger.Error("registration error", slog.String("msg", "get user failed"),
			slog.Int("status", http.StatusInternalServerError), slog.Any("error", err))
		return
	}

	if usecase.UserExists(getUser) {
		http.Error(w, "user already exists", http.StatusInternalServerError) //сервер не может обрабатывать запросы
		s.logger.Error("registration error", slog.String("msg", "user already exists"),
			slog.Int("status", http.StatusInternalServerError), slog.Any("error", err))
		return
	}
	err = s.u.RegisterUser(&user)
	if err != nil {
		http.Error(w, "registration is failed", http.StatusInternalServerError) //сервер не может обрабатывать запросы
		s.logger.Error("registration error", slog.String("msg", "registration is failed"),
			slog.Int("status", http.StatusInternalServerError), slog.Any("error", err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	s.logger.Info("registration successful", slog.Int("status", http.StatusCreated))
	//возращаем secret
	json.NewEncoder(w).Encode(user.Secret)
}

//func (s *Server) APIGetHandler(w http.ResponseWriter, r *http.Request) {}
//func (s *Server) APIConvertHandler(w http.ResponseWriter, r *http.Request) {}
//func (s *Server) APIHistoryHandler(w http.ResponseWriter, r *http.Request) {}
