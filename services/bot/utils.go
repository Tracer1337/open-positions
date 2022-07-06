package main

import (
	"net/http"
	"time"
)

func getWithRetries(url string, retries int) (*http.Response, error) {
	resp, err := http.Get(url)
	if !isSuccess(resp, err) && retries > 0 {
		if err == nil {
			resp.Body.Close()
		}
		return getWithRetries(url, retries-1)
	}
	return resp, err
}

func isSuccess(resp *http.Response, err error) bool {
	return err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300
}

func waitFor(fn func() bool, ms int, retries int) {
	i := 0
	for !fn() && i < retries {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		i++
	}
}
