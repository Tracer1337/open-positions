package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

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
		render(resp, filepath.Join(path, "README.md"))
	})
}

func render(resp *CompanyResponse, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error creating file %s\n", path)
	}
	defer file.Close()

	tmpl := template.Must(template.ParseFiles("template.md"))
	tmpl.Execute(file, resp.Data)
}
