package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
)

func main() {
	fmt.Println("Staring backend server...")

	db.Connect()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// logger, _ := zap.NewProduction()
	// defer logger.Sync()
	// sugar := logger.Sugar()

	// sugar.Infow("Test logging",
	// 	"url", "googli.com",
	// 	"num", 3,
	// 	"tries", 42.0)

	// logger.Error("Trying out logger",
	// 	zap.String("message", "meow"))

	recievedSignal := <-sigs

	fmt.Println("Recieved signal:", recievedSignal, "\nShutting down server...")
	db.Disconnect()
}
