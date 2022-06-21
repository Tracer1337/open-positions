package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type CompanyResponse struct {
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
	Meta struct {
		Pagination struct {
			Page      int `json:"page"`
			PageSize  int `json:"pageSize"`
			PageCount int `json:"pageCount"`
			Total     int `json:"total"`
		} `json:"pagination"`
	} `json:"meta"`
}

func FetchCompanies() (*CompanyResponse, error) {
	url := os.Getenv("STRAPI_URL") + "/api/companies?pagination[pageSize]=100"
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+os.Getenv("STRAPI_TOKEN"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CompanyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
