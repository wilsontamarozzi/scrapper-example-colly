package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

type pageInfo struct {
	StatusCode int
	Market     []Kripto
}

type Kripto struct {
	Name              string
	Symbol            string
	MarketCap         string
	Price             string
	CirculatingSupply string
	Volume24h         string
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := colly.NewCollector()

	p := &pageInfo{Market: []Kripto{}}

	c.OnHTML("#currencies-all tbody tr", func(e *colly.HTMLElement) {
		var coin = Kripto{
			Name:              e.ChildText(".currency-name-container"),
			Symbol:            e.ChildText(".col-symbol"),
			MarketCap:         e.ChildAttr(".market-cap", "data-usd"),
			Price:             e.ChildAttr("a.price", "data-usd"),
			CirculatingSupply: e.ChildAttr("td.circulating-supply span", "data-supply"),
			Volume24h:         e.ChildAttr("a.volume", "data-usd"),
		}

		p.Market = append(p.Market, coin)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	b, err := json.Marshal(p)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	addr := ":7171"

	http.HandleFunc("/", handler)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
