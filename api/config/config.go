package config

import (
	logger "github.com/horoshi10v/loggerNIX/v4"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

//type Config struct {
//	Port       string `toml:"port"`
//	LoggerPath string `toml:"logger_path"`
//}
//
//func NewConfig() *Config {
//	return &Config{
//		Port:       ":8080",
//		LoggerPath: "./logs/server.log",
//	}
//}

type Configer struct {
	Port                   string
	AccessSecret           string
	RefreshSecret          string
	AccessLifetimeMinutes  int
	RefreshLifetimeMinutes int
	Logger                 logger.Logger
	Driver                 string
	DataSourceName         string
}

func NewConfigs(l logger.Logger) *Configer {
	err := godotenv.Load("api/config/conf.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	accessLifetimeMinutes, _ := strconv.Atoi(os.Getenv("ACCESS_LIFE_TIME"))
	refreshLifetimeMinutes, _ := strconv.Atoi(os.Getenv("REFRESH_LIFE_TIME"))
	return &Configer{
		Port:                   os.Getenv("PORT"),
		AccessSecret:           os.Getenv("ACCESS_SECRET"),
		Logger:                 l,
		RefreshSecret:          os.Getenv("REFRESH_SECRET"),
		AccessLifetimeMinutes:  accessLifetimeMinutes,
		RefreshLifetimeMinutes: refreshLifetimeMinutes,
		Driver:                 os.Getenv("DRIVER"),
		DataSourceName:         os.Getenv("DATA_SOURCE_NAME"),
	}
}
