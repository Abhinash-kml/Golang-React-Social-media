package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Abhinash-kml/Golang-React-Social-media/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting backend server...")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Couldn't load .env file")
	}

	server := server.NewServer()
	server.Start()

	recievedSignal := <-sigs
	fmt.Println("Recieved signal:", recievedSignal, "\nShutting down server...")
	server.Stop()
}
