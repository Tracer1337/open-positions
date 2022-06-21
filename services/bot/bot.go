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
		exec("git", "commit", "--allow-empty", "-m", "\"test\"")
	})
}
