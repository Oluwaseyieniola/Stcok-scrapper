package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Stock struct {
	change, price, company string
}

func main() {
	ticker := []string{
		"MSFT",
		"IBM",
		"DIS",
		"ADBE",
		"AMZN",
		"WMT",
		"AL",
		"BCS",
		"BE",
		"AXP",
		"BA",
		"BEEP",
		"BEP",
		"BH",
		"CC",
		"CCL",
		"FTK",
	}

	stocks := []Stock{}
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("checking...", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Not right", err)
	})
	c.OnHTML("divBquote-header-info", func(c *colly.HTMLElement) {
		stock := Stock{}
		stock.company = c.ChildText("h1")
		fmt.Println("Company:", stock.company)
		stock.price = c.ChildText("fin-streamer[data-field='regularMarketPrice']")
		fmt.Println("price:", stock.price)
		stock.change = c.ChildText("fin-streamer[data-field='regularMarketChangePercent]")
		fmt.Println("change", stock.change)

		stocks = append(stocks, stock)
	})
	c.Wait()

	for _, t := range ticker {
		c.Visit("https://finance.yahoo.com/quote/ " + t + "/")
	}
	fmt.Println(stocks)

	file, err := os.Create("stocks.csv")

	if err != nil {
		log.Fatal("failed to create csv file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	headers := []string{
		"company",
		"price",
		"change",
	}
	writer.Write(headers)
	for _, stock := range stocks {
		record := []string{
			stock.company,
			stock.change,
			stock.price,
		}
		writer.Write(record)
	}

	defer writer.Flush()

}
