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

func GetDatabaseURL() string {
	url, err := checkEnv("DB_URL")
	if err != nil {
		url = "postgres://postgres:example@127.0.0.1:5432/dev_db?sslmode=disable"
	}

	return url
}

func GetTestDatabaseURL() string {
	url := os.Getenv("TEST_DB_URL")
	if url == "" {
		url = "postgres://postgres:example@127.0.0.1:5433/test_db?sslmode=disable"
	}

	return url
}

func GetServerBindAddress() string {
	host, err := checkEnv("SERVER_HOST")
	Must(err)
	port, err := checkEnv("SERVER_PORT")
	Must(err)

	return fmt.Sprintf("%s:%s", host, port)
}

func GetServerReadTimeout() string {
	v, _ := checkEnv("SERVER_READ_TIMEOUT")
	if v == "" {
		return "60"
	}

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
