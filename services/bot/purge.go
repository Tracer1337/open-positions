package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"open-positions/bot/api"
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
	resp, err := api.FetchCompanies()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wg.Add(len(resp.Data))

	for _, company := range resp.Data {
		go func(company api.Company) {
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

func checkCompany(company api.Company) bool {
	isValid := true

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		resp, err := http.Get(company.Attributes.WebsiteUrl)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
			isValid = false
		}
		wg.Done()
	}()

	go func() {
		resp, err := http.Get(company.Attributes.OpenPositionsUrl)
		if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
			isValid = false
		}
		wg.Done()
	}()

	wg.Wait()

	return isValid
}

func invalidateCompany(company api.Company) {
	json := "{ \"data\": { \"publishedAt\": null } }"

	_, err := api.FetchAPI("PUT", "/companies/"+fmt.Sprint(company.Id), api.RequestOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: bytes.NewBuffer([]byte(json)),
	})

	if err != nil {
		panic(err)
	}
}
