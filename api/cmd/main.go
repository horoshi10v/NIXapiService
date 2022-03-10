package main

import (
	"NIXSwag/api/config"
	"NIXSwag/api/internal/database"
	"NIXSwag/api/server"
	"database/sql"
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	logger "github.com/horoshi10v/loggerNIX/v4"
	"log"
)

func main() {
	//config
	myLog := logger.NewLogger("api/logs/logs.log")
	configs := config.NewConfigs(myLog)

	//database
	conn, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	//database.DeleteTables(conn)
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	////parser
	//err = p.Parser(conn, err)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//server
	s := server.New(configs, conn)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	myLog.Info("START")

}
