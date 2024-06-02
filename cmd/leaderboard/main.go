package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"

	memory_cache "github.com/nero2009/toptraders"
)

type TradingInfo struct {
	StockSymbol string
	Price       float64
	LastUpdated int64
	TraderName  string
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	http.HandleFunc("/leaderboard", leaderboardHandler)
	http.HandleFunc("/leaderboard/", leaderboardBySymbolHandler)
	http.ListenAndServe(":9091", nil)

	fmt.Print("Running leaderboard on port 9091")

}

func leaderboardBySymbolHandler(w http.ResponseWriter, r *http.Request) {
	var supportSymbols = []string{"BTC", "ETH", "LTC", "XRP", "BCH", "EOS", "XLM", "ADA", "TRX", "XMR"}
	var tradingInfo []TradingInfo
	path := r.URL.Path
	symbol := path[len("/leaderboard/"):]

	if symbol == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !contains(supportSymbols, symbol) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Symbol not supported"))
		return
	}

	//check if symbol is in cache
	cachedLeaders, _ := memory_cache.MemoryCache.Get(symbol)

	if cachedLeaders != nil {
		fmt.Printf("Cache hit, fetching leaderBoard for symbol %s from cache\n", symbol)
		err := json.Unmarshal(cachedLeaders, &tradingInfo)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tradingInfo)
		return
	}

	fmt.Printf("Cache miss, fetching leaderBoard for symbol %s from api\n", symbol)

	tradeInfo, err := getTopTraders()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var filteredTradingInfo []TradingInfo

	for _, trade := range tradeInfo {
		if trade.StockSymbol == symbol {
			filteredTradingInfo = append(filteredTradingInfo, trade)
		}
	}

	sort.Slice(filteredTradingInfo, func(i, j int) bool {
		return filteredTradingInfo[j].Price < filteredTradingInfo[i].Price
	})

	top10Trades := filteredTradingInfo[0:10]
	byteValue, err := json.Marshal(top10Trades)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//cache symbol
	memory_cache.MemoryCache.Set(symbol, byteValue)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(top10Trades)

}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var tradingInfo []TradingInfo

	cachedLeaderBoard, err := memory_cache.MemoryCache.Get("leaderboard")

	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "not found") {
			fmt.Println("Cache miss, fetching leaderBoard from api")

			info, err := getTopTraders()

			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			sort.Slice(info, func(i, j int) bool {
				return info[j].Price < info[i].Price
			})

			cachedInfo, err := json.Marshal(info)

			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			memory_cache.MemoryCache.Set(
				"leaderboard", cachedInfo)

			// cache the top 10 traders

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(info)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cachedLeaderBoard != nil {
		fmt.Println("Cache hit, fetching leaderBoard from cache")
		err = json.Unmarshal(cachedLeaderBoard, &tradingInfo)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tradingInfo)
		return
	}

}

func getTopTraders() ([]TradingInfo, error) {
	var tradingInfo []TradingInfo
	var err error
	url := os.Getenv("TRADING_INFO_URL")
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error connecting to external server")
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var errd = decoder.Decode(&tradingInfo)
	if errd != nil {
		fmt.Println(errd)
		return nil, fmt.Errorf("error decode response from trading site")
	}
	return tradingInfo, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
