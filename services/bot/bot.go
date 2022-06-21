package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file")
	}

	WithGitRepo(func(exec func(string, ...string), path string) {
		resp, err := FetchCompanies()
		if err != nil {
			log.Fatal(err)
		}
		for _, company := range resp.Data {
			log.Println(company.Attributes)
		}
	})
}
