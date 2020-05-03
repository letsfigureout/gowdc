package handlers

import (
	"fmt"
	"gowdc/internal/data"
	"log"
	"net/http"
	"strconv"

	"encoding/json"
)

type GoWDC struct {
	Log *log.Logger
}

func (g *GoWDC) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var payload interface{}
		var days, status int

		qs := r.URL.Query()
		numDays, OK := qs["days"]

		// default number of days
		days = 365
		if OK {
			var err error
			days, err = strconv.Atoi(numDays[0])
			if err != nil {
				payload = fmt.Sprintf("error parsing days : %v", err)
			}
		}

		ticker, OK := qs["ticker"]
		if !OK {
			g.Log.Println("error: missing ticker")
			payload = []string{"missing ticker"}
			status = http.StatusBadRequest
		} else {

			g.Log.Printf("Tickers: %v", ticker)
			datasource := data.NewYahooBackend(
				"https://query1.finance.yahoo.com/v7/finance/download/###?period1=###&period2=###&interval=1d&events=history",
			)

			stockData, err := datasource.GetStockHistory(ticker, days)
			if err != nil {
				payload = fmt.Sprintf("error getting stock data : %v", err)
				status = http.StatusInternalServerError
			}
			payload = stockData
		}

		if status != 0 {
			w.WriteHeader(status)
		}

		json.NewEncoder(w).Encode(payload)
	}
}
