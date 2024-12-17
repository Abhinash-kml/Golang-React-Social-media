package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Abhinash-kml/Golang-React-Social-media/pkg/db"
)

func main() {
	fmt.Println("Hello World")

	db.Connect()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	// Shutdown()
}
