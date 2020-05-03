package data

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Price of a stock
type Price struct {
	Date     time.Time `json:"Date"`
	Open     float64   `json:"Open"`
	High     float64   `json:"High"`
	Low      float64   `json:"Low"`
	Close    float64   `json:"Close"`
	AdjClose float64   `json:"AdjClose"`
	Volume   int64     `json:"Volume"`
}

// Stock details
type Stock struct {
	Prices  []Price `json:"PriceHistory"`
	Ticker  string  `json:"Ticker"`
	DataURL string  `json:"-"`
}

type StockList struct {
	Stock map[string]Stock `json:"Stock"`
}

// YahooData endpoint
type YahooData struct {
	QuoteURL string
}

func (y *YahooData) GetStockHistory(tickers []string, days int) (StockList, error) {

	// The stock history url is in the following format:
	// https://query1.finance.yahoo.com/v7/finance/download/###?period1=1556706946&period2=1588329346&interval=1d&events=history
	// ### is the ticker identifier in uppercase

	now := time.Now().UTC()
	end := now.Unix()
	start := now.AddDate(0, 0, -days).Unix()

	var sl StockList
	sl.Stock = make(map[string]Stock, len(tickers))

	for _, ticker := range tickers {

		ticker = strings.ToUpper(ticker)

		tickerQuote := y.buildQuoteURL(ticker, start, end)
		log.Printf("Ticker Quote URL: %v", tickerQuote)

		prices, err := y.downloadCSVData(tickerQuote)
		if err != nil {
			return StockList{}, fmt.Errorf("error getting prices: %v\n", err)
		}

		stockHistory := Stock{
			Prices:  prices,
			Ticker:  ticker,
			DataURL: tickerQuote,
		}

		sl.Stock[ticker] = stockHistory
	}

	return sl, nil
}

// NewYahooBackend factory function for Yahoo Finance data backend
func NewYahooBackend(quoteURL string) YahooData {

	yd := YahooData{
		QuoteURL: quoteURL,
	}

	return yd
}

func (y *YahooData) buildQuoteURL(ticker string, start, end int64) string {

	// Replace ticker
	q := strings.Replace(y.QuoteURL, "###", ticker, 1)

	// Parse string URL into struct so we can replace the query components
	tickerURL, err := url.ParseRequestURI(q)
	if err != nil {
		log.Printf("error parsing\n")
	}

	m, err := url.ParseQuery(tickerURL.RawQuery)
	if err != nil {
		log.Printf("error parsing\n")
	}

	// Replace start and end dates
	m["period1"] = []string{strconv.Itoa(int(start))}
	m["period2"] = []string{strconv.Itoa(int(end))}

	tickerURL.RawQuery = m.Encode()
	return tickerURL.String()
}

func (y *YahooData) downloadCSVData(csvURL string) ([]Price, error) {

	log.Printf("downloading csv data from %v", csvURL)
	r, err := http.Get(csvURL)
	if err != nil || r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to download csv data: [%v] %v", r.StatusCode, err)
	}
	defer r.Body.Close()

	csvData := csv.NewReader(r.Body)

	// Header, we aren't going to use it here so discard it
	_, err = csvData.Read()
	if err != nil {
		return nil, fmt.Errorf("error parsing csv data : %v", err)
	}

	records, err := csvData.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing csv data : %v", err)
	}

	var prices []Price
	for _, row := range records {

		rDate, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing csv data : %v", err)
		}

		rOpen, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing csv data : %v", err)
		}

		rHigh, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing csv data : %v", err)
		}

		rLow, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing csv data : %v", err)
		}

		rClose, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing csv data : %v", err)
		}

		rAdjClose, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing csv data : %v", err)
		}

		rVolume, err := strconv.ParseInt(row[6], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing csv data : %v", err)
		}

		p := Price{
			Date:     rDate,
			Open:     rOpen,
			High:     rHigh,
			Low:      rLow,
			Close:    rClose,
			AdjClose: rAdjClose,
			Volume:   rVolume,
		}

		prices = append(prices, p)
	}

	return prices, nil
}
