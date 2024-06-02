package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type TradingInfo struct {
	StockSymbol string
	Price       float64
	LastUpdated int64
	TraderName  string
}

func main() {
	http.HandleFunc("/tradinginfo", tradinginfoHandler)
	http.ListenAndServe(":9090", nil)
}

func tradinginfoHandler(w http.ResponseWriter, r *http.Request) {
	tradingInfo := getTradingInfo()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tradingInfo)
}

// generate a new slice of TradingInfo structs anytime this function is called

func getTradingInfo() []TradingInfo {
	var tradingInfo = make([]TradingInfo, 0)
	var cryptoSymbols = []string{"BTC", "ETH", "LTC", "XRP", "BCH", "EOS", "XLM", "ADA", "TRX", "XMR"}

	for i := 0; i < 1000; i++ {
		tradingInfo = append(tradingInfo, TradingInfo{
			StockSymbol: cryptoSymbols[gofakeit.Number(0, len(cryptoSymbols)-1)],
			Price:       gofakeit.Price(1, 1000),
			LastUpdated: time.Now().Unix(),
			TraderName:  gofakeit.Name(),
		})
	}

	return tradingInfo
}
