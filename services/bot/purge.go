package main

import (
	"fmt"
	"log"
	"open-positions/bot/api"
	"sync"

	"github.com/robfig/cron/v3"
)

func schedulePurge() {
	c := cron.New()

	c.AddFunc("@hourly", func() {
		log.Println("Purge Start")
		runPurge()
		log.Println("Purge Done")
	})

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
	if company.Attributes.OpenPositionsCount == 0 {
		return false
	}

	isValid := true

	var wg sync.WaitGroup
	var m sync.Mutex

	wg.Add(2)

	go func() {
		resp, err := getWithRetries(company.Attributes.WebsiteUrl, 3)
		result := isSuccess(resp, err)
		if err == nil {
			defer resp.Body.Close()
		}
		m.Lock()
		isValid = result
		m.Unlock()
		wg.Done()
	}()

	go func() {
		resp, err := getWithRetries(company.Attributes.OpenPositionsUrl, 3)
		result := isSuccess(resp, err)
		if err == nil {
			defer resp.Body.Close()
		}
		m.Lock()
		isValid = result
		m.Unlock()
		wg.Done()
	}()

	wg.Wait()

	return isValid
}

func invalidateCompany(company api.Company) {
	_, err := api.FetchAPI("PUT", "/companies/"+fmt.Sprint(company.Id), api.RequestOptions{
		Body: "{ \"data\": { \"publishedAt\": null } }",
	})
	if err != nil {
		panic(err)
	}
}
