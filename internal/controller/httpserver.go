package controller

import (
	"crypto-project/internal/usecase"
	"crypto-project/pkg/logger"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
)

type Server struct { //controller
	router    *mux.Router
	u         *usecase.Usecase
	logger    *logger.Logger
	secretKey []byte
}

func New(u *usecase.Usecase, logger *logger.Logger) *Server {
	s := &Server{router: mux.NewRouter(),
		u:         u,
		logger:    logger,
		secretKey: []byte("your_secret_key"),
	}
	s.router.HandleFunc("/home", s.HomeHandler).Methods("GET")
	s.router.HandleFunc("/info", s.InfoHandler).Methods("GET")

	s.router.HandleFunc("/login", s.LoginHandler).Methods("POST")
	s.router.HandleFunc("/register", s.RegisterHandler).Methods("POST")

	s.router.HandleFunc("/api/get", s.APIGetHandler).Methods("GET")
	s.router.HandleFunc("/api/convert", s.APIConvertHandler).Methods("POST")
	s.router.HandleFunc("/api/history", s.APIHistoryHandler).Methods("GET")

	s.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return s
}

func (s *Server) Run(port string) {
	s.logger.Info("Сервер запущен на http://127.0.0.1:" + port)
	if err := http.ListenAndServe("localhost:"+port, s.router); err != nil {
		s.logger.Error("fatal error", slog.Int("status", http.StatusBadGateway))
		return
	}
}
