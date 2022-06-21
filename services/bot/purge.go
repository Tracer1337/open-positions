package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/robfig/cron/v3"
)

func schedulePurge() {
	c := cron.New()

	c.AddFunc("@daily", runPurge)

	c.Start()

	select {}
}

func runPurge() {
	resp, err := fetchCompanies()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wg.Add(len(resp.Data))

	for _, company := range resp.Data {
		go func(company Company) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Error checking " + company.Attributes.Name)
				}
			}()
			if !checkCompany(company) {
				invalidateCompany(company)
			}
			wg.Done()
		}(company)
	}

	wg.Wait()
}

func checkCompany(company Company) bool {
	resp, err := http.Get(company.Attributes.WebsiteUrl)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false
	}

	resp, err = http.Get(company.Attributes.OpenPositionsUrl)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false
	}

	return true
}

func invalidateCompany(company Company) {
	log.Println("Invalidate " + company.Attributes.Name)

	json := "{ \"data\": { \"publishedAt\": null } }"

	body, err := fetchAPI("PUT", "/companies/"+fmt.Sprint(company.Id), RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: bytes.NewBuffer([]byte(json)),
	})

	if err != nil {
		panic(err)
	}

	fmt.Print(string(body))
}
