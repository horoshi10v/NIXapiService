package main

import (
	config2 "NIXSwag/api/config"
	"NIXSwag/api/internal/database"
	p "NIXSwag/api/internal/parser"
	"NIXSwag/api/server"
	"database/sql"
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	logger "github.com/horoshi10v/loggerNIX/v4"
	"log"
)

//
//var (
//	configPath string
//)
//
//func init() {
//	flag.StringVar(&configPath,
//		"config-path",
//		"api/config/conf.env",
//		"path to config file")
//}

func main() {
	//config
	//flag.Parse()
	myLog := logger.NewLogger("api/logs/logs.log")
	config := config2.NewConfigs(myLog)
	//_, err := toml.DecodeFile(configPath, config)
	//if err != nil {
	//	log.Fatal(err)
	//}

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
	//server
	s := server.New(config, conn)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	myLog.Info("START")
	//parser
	err = p.Parser(conn, err)
	if err != nil {
		return
	}

}
