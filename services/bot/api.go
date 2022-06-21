package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/go-querystring/query"
)

type RequestOptions struct {
	Pagination RequestPagination `url:"pagination"`
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

func fetchAPI(path string, opts RequestOptions) ([]byte, error) {
	v, err := query.Values(opts)
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("%s/api%s?%s", os.Getenv("STRAPI_URL"), path, v.Encode())

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("STRAPI_TOKEN"))

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

type CompaniesResponse struct {
	Data []struct {
		Id         int `json:"id"`
		Attributes struct {
			Name               string `json:"name"`
			WebsiteUrl         string `json:"website_url"`
			OpenPositionsCount int    `json:"open_positions_count"`
			OpenPositionsUrl   string `json:"open_positions_url"`
			EmployeesCount     int    `json:"employees_count"`
			Revenue            string `json:"revenue"`
		} `json:"attributes"`
	} `json:"data"`
	Meta ResponseMeta `json:"meta"`
}

func FetchCompanies() (*CompaniesResponse, error) {
	body, err := fetchAPI("/companies", RequestOptions{
		Pagination: RequestPagination{
			PageSize: 100,
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
