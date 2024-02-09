package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func fetchDocument(url string) (*html.Node, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Parse the page
	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get("https://www.investing.com/etfs/spdr-s-p-500")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		fmt.Printf("Title: %s\n", title)
	})
	document, err := fetchDocument("https://www.investing.com/equities/tesla-motors")
	if err != nil {
		log.Fatal(err)
	}
	xpath := `/html/body/div[1]/div[2]/div[2]/div[2]/div[1]/div[1]/div[1]/div[1]/h1[1]`
	node := htmlquery.FindOne(document, xpath)

	priceXpath := `//*[@id="__next"]/div[2]/div[2]/div[2]/div[1]/div[1]/div[3]/div[1]/div[3]/div[2]/span[1]`
	fmt.Println("Found element index:", htmlquery.InnerText(node))
	node = htmlquery.FindOne(document, priceXpath)
	fmt.Println("Found element price:", htmlquery.InnerText(node))

}
func main() {
	ExampleScrape()
}
