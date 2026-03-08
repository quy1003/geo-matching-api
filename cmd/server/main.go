package main

import (
	"log"

	"github.com/quy1003/geo-matching-api/internal/app"
)

func main() {
	server := app.NewServer()
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
