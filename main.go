package main

import (
	"fmt"

	"github.com/invest-scraping/config"
	"github.com/invest-scraping/scrape/goquery"
)

func main() {
	cfg := config.Load("./config/env")

	for _, m := range cfg.Monitors {

		fmt.Printf("Starting monitor for %s\n", m.Symbol)
		// MontorScrape(m.Endpoint, m.EndpointExt, m.NameXpath, m.PriceXpath)
		scrap := goquery.NewScraper(m.Endpoint+m.EndpointExt, m.PriceXpath)
		resp, err := scrap.RunQuery()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println("Monitor: ", m.Name, "Price: ", resp)
		// go blockIndexer.Start()
	}
	// ExampleScrape()
}
