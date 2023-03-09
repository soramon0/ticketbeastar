package utils

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func StartServer(a *fiber.App, l *log.Logger) {
	if GetStageStatus() == "dev" {
		listenAndServe(a, l)
	} else {
		go listenAndServe(a, l)
		startServerWithGracefulShutdown(a, l)
	}
}

func listenAndServe(a *fiber.App, l *log.Logger) {
	if err := a.Listen(GetServerBindAddress()); err != nil {
		l.Fatalf("Oops... Server is not running! Reason: %v\n", err)
	}
}

func startServerWithGracefulShutdown(a *fiber.App, l *log.Logger) {
	// trap interupt, sigterm or sighub and gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	// Block until a signal is received.
	sig := <-sigChan
	l.Printf("Recieved %s, graceful shutdown...\n", sig)

	// gracefully shutdown the server
	if err := a.Shutdown(); err != nil {
		l.Fatal(err)
	}
}
