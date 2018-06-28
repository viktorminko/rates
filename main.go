package main

import (
	"log"
	"time"
	"flag"
	"net/url"
	"strings"
	"net/http"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

const defDuration = 5 * time.Second
const defUpdateURL = "https://wex.nz/api/3/ticker/"
const defPairsURL = "https://wex.nz/api/3/info/"
const defEndpointURL = "/"
const defTimeFrame = 10 * time.Minute
const defPort = 3512

type Config struct {
	UpdateDuration time.Duration
	UpdateURL      url.URL
	EndpointURL    url.URL
	TimeFrame      time.Duration
	Pairs          []string
	Port           int
}

type Rate struct {
	Pair  string
	Value float64
}

type RateTime struct {
	Rate Rate
	Time time.Time
}

type Averages = sync.Map

func RunUpdater(url url.URL, duration time.Duration, pairs []string) (<-chan RateTime) {

	ch := make(chan RateTime)

	go func() () {
		defer close(ch)
		t := time.NewTicker(duration)
		for {
			for i := range pairs {
				go APIGetStatsForPair(ch, pairs[i], url)
			}
			<-t.C
		}
	}()

	return ch
}

func RunCalculator(chIn <-chan RateTime, chSignal <-chan struct{}, timeSpan time.Duration) (<-chan Averages) {

	curAverage := Averages{}
	queue := map[string][]float64{}

	timeStarted := time.Now()

	ch := make(chan Averages)

	go func() () {
		defer close(ch)
		for {
			select {
			case rateTime := <-chIn:

				r := rateTime.Rate

				UpdateAverages(
					r.Pair,
					r.Value,
					&queue,
					&curAverage,
					rateTime.Time.Sub(timeStarted) <= timeSpan,
				)

			case <-chSignal:
				ch <- curAverage
			}

		}
	}()

	return ch
}

func Handler(chSignal chan<- struct{}, chOut <-chan Averages, w http.ResponseWriter) () {
	chSignal <- struct{}{}

	averages := <-chOut

	strAverages := map[string]string{}

	averages.Range(func(k, v interface{}) bool {
		strAverages[k.(string)] = fmt.Sprintf("%.3f", v)
		return true
	})

	js, err := json.MarshalIndent(strAverages, "", "   ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func InitConfig() (*Config, error) {
	duration := flag.Duration("update_duration", defDuration, "how often rates are updated")
	pairsURLStr := flag.String("get_pairs_url", defPairsURL, "URL to get pairs information")
	updateURLStr := flag.String("update_url", defUpdateURL, "URL to update rates for pair")
	endpointURLStr := flag.String("api_endpoint", defEndpointURL, "API endpoint to request average rates")
	timeFrame := flag.Duration("timeframe", defTimeFrame, "time interval to use for moving averages")
	pairsStr := flag.String("pairs", "", "comma separated currency pairs : btc_usd, eth_eur, xrp_btc")
	port := flag.Int("port", defPort, "port for connections")

	flag.Parse()

	pairsURL, err := url.Parse(*pairsURLStr)
	if err != nil {
		pairsURL, _ = url.Parse(defPairsURL)
	}

	updateURL, err := url.Parse(*updateURLStr)
	if err != nil {
		updateURL, _ = url.Parse(defUpdateURL)
	}

	endpointURL, err := url.Parse(*endpointURLStr)
	if err != nil {
		endpointURL, _ = url.Parse(defEndpointURL)
	}

	var pairs []string
	if len(*pairsStr) > 0 {
		pairs = strings.Split(*pairsStr, ",")
	}

	//If no pairs provided, we get all from API
	if 0 == len(pairs) {
		pairs, err = APIGetAllPairs(*pairsURL)
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		*duration,
		*updateURL,
		*endpointURL,
		*timeFrame,
		pairs,
		*port,
	}, nil
}

func main() {

	log.Println("Service started")

	config, err := InitConfig()
	if err != nil {
		log.Fatalf("unable to load currency pairs, error: %v", err)
	}

	confStr, err := json.MarshalIndent(config, "", "   ")
	if err != nil {
		log.Printf("error while unmarshalling config: %v", err)
	}
	log.Printf("Config initialized: %s", confStr)

	ch := RunUpdater(config.UpdateURL, config.UpdateDuration, config.Pairs)

	chSignal := make(chan struct{})
	defer close(chSignal)

	chOut := RunCalculator(ch, chSignal, config.TimeFrame)

	http.HandleFunc(config.EndpointURL.Path, func(w http.ResponseWriter, r *http.Request) {
		Handler(chSignal, chOut, w)
	})
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}
