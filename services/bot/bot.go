package main

//go:generate go install github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates

import (
	"fmt"
	"log"
	"open-positions/bot/api"
	"open-positions/bot/templates"
	"os"
	"path/filepath"
	"sync"
	"text/template"

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

	go scheduleScrape()

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

	companiesResp, err := api.FetchCompanies()
	if err != nil {
		panic(err)
	}

	readmeResp, err := api.FetchReadme()
	if err != nil {
		panic(err)
	}

	render(companiesResp.Data, readmeResp.Data, filepath.Join(path, "README.md"))

	if diff := exec("git", "diff", "README.md"); len(diff) == 0 {
		log.Println("Empty diff: Skip update")
		return
	}

	exec("git", "add", "README.md")
	exec("git", "commit", "-m", "chore: update readme")

	if _, isDryRun := os.LookupEnv("DRY_RUN"); isDryRun {
		log.Println("Dry-Run: Skip git push")
		return
	}

	exec("git", "push")
}

func render(companies []api.Company, readme api.Readme, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("Error creating file %s\n", path))
	}
	defer file.Close()

	tmpl, err := template.New("readme").Parse(readme.Attributes.Template)
	if err != nil {
		panic("Error parsing readme template")
	}

	table := templates.ReadmeTable(companies)
	err = tmpl.Execute(file, map[string]string{
		"Table": table,
	})
	if err != nil {
		panic("Error rendering table")
	}
}
