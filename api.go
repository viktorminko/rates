package main

import (
	"net/url"
	"net/http"
	"fmt"
	"encoding/json"
	"time"
)

func APIGetAllPairs(url url.URL) ([]string, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("invalid response status code, expected %v got $v", http.StatusOK, resp.Status)
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body, error: %v, body: %v", err, resp.Body)
	}

	pairsData, ok := data["pairs"]
	if !ok {
		return nil, fmt.Errorf("no pairs returned")
	}

	var pairs []string
	for k := range pairsData.(map[string]interface{}) {
		pairs = append(pairs, k)
	}

	return pairs, nil
}

func APIGetStatsForPair(c1 chan RateTime, pair string, url url.URL) (error) {
	resp, err := http.Get(url.String()+pair)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if http.StatusOK != resp.StatusCode {
		return fmt.Errorf("invalid response status code, expected %v got $v", http.StatusOK, resp.Status)
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("error decoding response body, error: %v, body: %v", err, resp.Body)
	}

	res, ok := data[pair]
	if !ok {
		return fmt.Errorf("no rate returned")
	}

	lastRate, ok := res.(map[string]interface{})["last"]
	if !ok {
		return fmt.Errorf("no last returned")
	}

	c1 <- RateTime{Rate{pair, lastRate.(float64)}, time.Now()}
	return nil
}