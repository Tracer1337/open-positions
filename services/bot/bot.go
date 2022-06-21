package main

//go:generate go get -d github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates

import (
	"fmt"
	"log"
	"open-positions/bot/api"
	"open-positions/bot/templates"
	"os"
	"path/filepath"
	"sync"

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

	go schedulePurge()

	r := gin.Default()

	var lock sync.Mutex

	r.POST("/update", func(c *gin.Context) {
		lock.Lock()
		updateReadme()
		lock.Unlock()
		c.Status(200)
	})

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}

func updateReadme() {
	exec, path := initGitRepo()

	resp, err := api.FetchCompanies()
	if err != nil {
		panic(err)
	}

	render(resp, filepath.Join(path, "README.md"))

	exec("git", "add", "README.md")
	exec("git", "commit", "--allow-empty", "-m", "chore: update readme")

	if _, isDryRun := os.LookupEnv("DRY_RUN"); isDryRun {
		log.Println("Dry-Run: Skip git push")
		return
	}

	exec("git", "push")
}

func render(resp *api.CompaniesResponse, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("Error creating file %s\n", path))
	}
	defer file.Close()

	content := templates.Readme(resp.Data)
	file.WriteString(content)
}
