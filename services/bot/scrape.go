package main

import (
	"bytes"
	"context"
	"fmt"
	"open-positions/bot/api"
	"sync"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func runScrape() {
	companies, err := api.FetchCompanies()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wg.Add(len(companies.Data))

	for _, company := range companies.Data {
		go func(company api.Company) {
			count := scrapeOpenPositions(company)
			if count != company.Attributes.OpenPositionsCount {
				body := fmt.Sprintf("{ \"data\": { \"open_positions_count\": %d } }", count)
				api.FetchAPI("PUT", "/companies/"+fmt.Sprint(company.Id), api.RequestOptions{
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: bytes.NewBuffer([]byte(body)),
				})
			}
			wg.Done()
		}(company)
	}

	wg.Wait()
}

func scrapeOpenPositions(company api.Company) int {
	if company.Attributes.OpenPositionsSelector == "" {
		return company.Attributes.OpenPositionsCount
	}

	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), "ws://localhost:3000")
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
