package main

import (
	"context"
	"fmt"
	"log"
	"open-positions/bot/api"
	"os"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/robfig/cron/v3"
)

func scheduleScrape() {
	c := cron.New()

	c.AddFunc("@daily", func() {
		log.Println("Scrape Start")
		runScrape()
		log.Println("Scrape Done")
	})

	c.Start()

	select {}
}

func runScrape() {
	companies, err := api.FetchCompanies()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wg.Add(len(companies.Data))

	for _, company := range companies.Data {
		go func(company api.Company) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Error scraping", company.Attributes.Name, r)
				}
			}()

			defer wg.Done()

			count := scrapeOpenPositionsHtml(company)
			if count == 0 {
				count = scrapeOpenPositionsChrome(company)
			}

			if count != company.Attributes.OpenPositionsCount {
				log.Println("Update positions count", company.Attributes.Name, count)
				body := fmt.Sprintf("{ \"data\": { \"open_positions_count\": %d } }", count)
				api.FetchAPI("PUT", "/companies/"+fmt.Sprint(company.Id), api.RequestOptions{
					Body: body,
				})
			}
		}(company)
	}

	wg.Wait()
}

func scrapeOpenPositionsHtml(company api.Company) int {
	if company.Attributes.OpenPositionsSelector == "" {
		return company.Attributes.OpenPositionsCount
	}

	resp, err := getWithRetries(company.Attributes.OpenPositionsUrl, 3)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if !isSuccess(resp, err) {
		return 0
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	return doc.Find(company.Attributes.OpenPositionsSelector).Length()
}

func scrapeOpenPositionsChrome(company api.Company) int {
	if company.Attributes.OpenPositionsSelector == "" {
		return company.Attributes.OpenPositionsCount
	}

	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), os.Getenv("CHROME_URL"))
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	count := 0
	err := chromedp.Run(ctx,
		chromedp.Navigate(company.Attributes.OpenPositionsUrl),
		chromedp.ActionFunc(func(ctx context.Context) error {
			waitFor(func() bool {
				node, err := dom.GetDocument().Do(ctx)
				if err != nil {
					return false
				}

				nodes, err := dom.QuerySelectorAll(node.NodeID, company.Attributes.OpenPositionsSelector).Do(ctx)
				if err != nil {
					return false
				}

				count = len(nodes)

				return count > 0
			}, 200, 20)

			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	return count
}
