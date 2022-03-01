package server

import (
	"NIXSwag/api/config"
	"NIXSwag/api/internal/repositories"
	"NIXSwag/api/internal/repositories/provider"
	"NIXSwag/api/route"
	"database/sql"
	"fmt"
	logger "github.com/horoshi10v/loggerNIX/v4"
	"log"
	"net/http"
	"os"
)

type APIServer struct {
	config *config.Configer
	myLog  *logger.Logger
	conn   *sql.DB
}

func New(conf *config.Configer, db *sql.DB) *APIServer {
	return &APIServer{
		config: conf,
		conn:   db,
	}
}
func (s *APIServer) Start() error {
	conn := s.conn
	TX, err := conn.Begin()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	myLogger := s.config.Logger
	myLogger.Info("START SERVER")
	repositiriesProv := &provider.Provider{
		ProductRepository: repositories.NewProductRepo(conn, TX, s.myLog),
	}
	mux := http.NewServeMux()
	route.Route(*repositiriesProv, mux)
	log.Fatal(http.ListenAndServe(s.config.Port, mux))
	//s.myLog.Info("START SERVER")
	return err
}
