package config

import (
	"os"
)

type Config struct {
	ServerPort string
	BaseURL    string
	SecretKey  string
	UseMySQL   bool
	MySQLDSN   string
	StaticDir  string
}

var GlobalConfig = Config{
	ServerPort: "8080",
	BaseURL:    "http://127.0.0.1:8080",
	SecretKey:  "tiktok-secret-key-2026",
	UseMySQL:   false,
	MySQLDSN:   "root:123456@tcp(127.0.0.1:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local",
	StaticDir:  "./static",
}

func InitConfig() {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		GlobalConfig.ServerPort = port
	}
	if host := os.Getenv("BASE_URL"); host != "" {
		GlobalConfig.BaseURL = host
	}
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		GlobalConfig.MySQLDSN = dsn
		GlobalConfig.UseMySQL = true
	}
}
