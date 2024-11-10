package controller

import (
	"context"
	"crypto-project/internal/usecase"
	"crypto-project/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"time"
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
	s.router.Use(s.SetHeaders)
	s.router.HandleFunc("/home", s.HomeHandler).Methods("GET")
	s.router.HandleFunc("/info", s.InfoHandler).Methods("GET")

	s.router.HandleFunc("/login", s.LoginHandler).Methods("POST")
	s.router.HandleFunc("/register", s.RegisterHandler).Methods("POST")

	apiRouter := s.router.PathPrefix("/api").Subrouter()
	apiRouter.Use(s.checkToken)
	apiRouter.HandleFunc("/get", s.APIGetHandler).Methods("GET")
	apiRouter.HandleFunc("/convert", s.APIConvertHandler).Methods("POST")
	apiRouter.HandleFunc("/history", s.APIHistoryHandler).Methods("GET")

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

func (s *Server) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userID", claims.UserID)))
	})
}

func (s *Server) SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("Request-ID", id)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("date", time.Now().Format("02.01.2006 15:04:05 Monday"))
		next.ServeHTTP(w, r)
	})
}
