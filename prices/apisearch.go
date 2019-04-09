package prices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type CryptoCompareResponse struct {
	USD            float64 `json:"USD"`
	ConversionType struct {
		Type             string `json:"type"`
		ConversionSymbol string `json:"conversionSymbol"`
	} `json:"ConversionType"`
}

func PriceSearch(coin, searchdate string) (*PriceItem, error) {

	log.Info("Searching Online for Price of Coin", coin, searchdate)
	safeCoin := url.QueryEscape(strings.ToUpper(coin))
	safeTime := url.QueryEscape(searchdate)

	layout := "2006-01-02"

	t, err := time.Parse(layout, safeTime)
	if err != nil {
		log.Warn("Could not parse the date parameter ", err)
		return nil, err
	}

	timestamp := int32(t.Unix())
	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/dayAvg?fsym=%s&tsym=USD&UTCHourDiff=10&toTs=%d&api_key=%s", safeCoin, timestamp, apiKey)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Warn("Could not build the request", err)
		return nil, err
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Warn("Could not send the request", err)
		return nil, err
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var record CryptoCompareResponse

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Warn("Could not decode the JSON request", err)
		return nil, err
	}
	log.Info("Received Price: USD = ", record.USD)

	response := &PriceItem{
		Currency: safeCoin,
		Base:     "USD",
		Date:     t.Format("2006-01-02"),
		Amount:   record.USD,
	}

	return response, nil
}
