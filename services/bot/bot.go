package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file")
	}

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.POST("/update", func(c *gin.Context) {
		updateReadme()
		c.Status(200)
	})

	r.Run("localhost:8000")
}

func updateReadme() {
	exec, path := InitGitRepo()

	resp, err := FetchCompanies()
	if err != nil {
		panic(err)
	}

	render(resp, filepath.Join(path, "README.md"))

	exec("git", "add", "README.md")
	exec("git", "commit", "--allow-empty", "-m", "chore: update readme")
	exec("git", "push")
}

func render(resp *CompanyResponse, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("Error creating file %s\n", path))
	}
	defer file.Close()

	tmpl := template.Must(template.ParseFiles("template.md"))
	tmpl.Execute(file, resp.Data)
}
