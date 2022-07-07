package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/go-querystring/query"
)

type RequestOptions struct {
	Body    string
	Headers map[string]string
	Query   RequestQuery
}

type RequestQuery struct {
	Pagination RequestPagination `url:"pagination"`
	Sort       []string          `url:"sort"`
}

type RequestPagination struct {
	PageSize int `url:"pageSize"`
}

type ResponseMeta struct {
	Pagination struct {
		Page      int `json:"page"`
		PageSize  int `json:"pageSize"`
		PageCount int `json:"pageCount"`
		Total     int `json:"to"`
	} `json:"pagination"`
}

func FetchAPI(method string, path string, opts RequestOptions) ([]byte, error) {
	if method != "GET" {
		if _, isDryRun := os.LookupEnv("DRY_RUN"); isDryRun {
			log.Printf("Dry-Run: Skip %s %s\n", method, path)
			return []byte{}, nil
		}
	}

	v, err := query.Values(opts.Query)
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("%s/api%s?%s", os.Getenv("STRAPI_URL"), path, v.Encode())

	client := http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(opts.Body)))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("STRAPI_TOKEN"))

	if method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	for k, v := range opts.Headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

type Company struct {
	Id         int `json:"id"`
	Attributes struct {
		Name                  string `json:"name"`
		WebsiteUrl            string `json:"website_url"`
		OpenPositionsCount    int    `json:"open_positions_count"`
		OpenPositionsUrl      string `json:"open_positions_url"`
		OpenPositionsSelector string `json:"open_positions_selector"`
		EmployeesCount        int    `json:"employees_count"`
		ImageUrl              string `json:"image_url"`
	} `json:"attributes"`
}

type CompaniesResponse struct {
	Data []Company    `json:"data"`
	Meta ResponseMeta `json:"meta"`
}

func FetchCompanies() (*CompaniesResponse, error) {
	body, err := FetchAPI("GET", "/companies", RequestOptions{
		Query: RequestQuery{
			Pagination: RequestPagination{
				PageSize: 100,
			},
			Sort: []string{"employees_count:desc"},
		},
	})
	if err != nil {
		panic(err)
	}

	var result CompaniesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

type Readme struct {
	Id         int `json:"id"`
	Attributes struct {
		Template string `json:"template"`
	} `json:"attributes"`
}

type ReadmeResponse struct {
	Data Readme `json:"data"`
}

func FetchReadme() (*ReadmeResponse, error) {
	body, err := FetchAPI("GET", "/readme", RequestOptions{})
	if err != nil {
		panic(err)
	}

	var result ReadmeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
