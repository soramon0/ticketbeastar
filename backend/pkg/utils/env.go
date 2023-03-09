package utils

import (
	"fmt"
	"os"
)

func checkEnv(envName string) (string, error) {
	v := os.Getenv(envName)
	if v == "" {
		return "", fmt.Errorf(`env variable "%s" is not defined`, envName)
	}
	return v, nil
}

func GetDatabaseBindAdress() string {
	proto, err := checkEnv("DB_PORTOCAL")
	Must(err)
	user, err := checkEnv("DB_USER")
	Must(err)
	pass, err := checkEnv("DB_PASSWORD")
	Must(err)
	host, err := checkEnv("DB_HOST")
	Must(err)
	port, err := checkEnv("DB_PORT")
	Must(err)
	dbName := GetDatabaseName()

	// eg. mongodb://user:password@mongo:27017/
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", proto, user, pass, host, port, dbName)
}

func GetServerBindAddress() string {
	host, err := checkEnv("SERVER_HOST")
	Must(err)
	port, err := checkEnv("SERVER_PORT")
	Must(err)

	return fmt.Sprintf("%s:%s", host, port)
}

func GetServerReadTimeout() string {
	v, err := checkEnv("SERVER_READ_TIMEOUT")
	Must(err)

	return v
}

func GetDatabaseName() string {
	v, err := checkEnv("DB_NAME")
	Must(err)

	return v
}

func GetStageStatus() string {
	v := os.Getenv("STAGE_STATUS")
	if v == "" {
		return "dev"
	}
	return v
}
