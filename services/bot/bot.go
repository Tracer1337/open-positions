package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file")
	}
	CloneRepo()
}
