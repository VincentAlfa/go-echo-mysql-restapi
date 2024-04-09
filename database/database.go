package database

import (
	"os"
	"github.com/joho/godotenv"
)

func DbSourceName() (string, string, string, string, string) {

	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	return user, password, host, port , database
}