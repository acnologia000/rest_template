package api

import (
	"database/sql"
	"net/http"
	"rest_template/internal/api/middleware"
	"rest_template/internal/config"
	"rest_template/shared/logger"
)

type Server struct {
	config *config.Config
	logger *logger.Logger
	db     *sql.DB
	router *http.ServeMux
}

func NewServer(cfg *config.Config, logger *logger.Logger, db *sql.DB) *Server {
	s := &Server{
		config: cfg,
		logger: logger,
		db:     db,
		router: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.router.Handle("/api/users", middleware.Chain(
		http.HandlerFunc(s.handleUsers),
		middleware.Logging(s.logger),
	))
	// Add more routes here
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	// Implement user handling logic
}
