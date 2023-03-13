package utils

import (
	"log"
	"os"
)

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "[soramon0/ticketBeast] ", log.LstdFlags)
}
