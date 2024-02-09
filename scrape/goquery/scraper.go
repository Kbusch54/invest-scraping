package goquery

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Scraper struct {
	// some fields
	document    *html.Node
	url         string
	path        string
	returnValue string
}

var (
	ErrNotFound = errors.New("Not found")
)

func NewScraper(url, path string) *Scraper {
	return &Scraper{
		url:  url,
		path: path,
	}
}

func (s *Scraper) GetDocument() error {
	res, err := http.Get(s.url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Parse the page
	doc, err := htmlquery.Parse(res.Body)
	if err != nil {
		return err
	}
	s.document = doc

	return nil
}

func (s *Scraper) FindElements() error {
	node := htmlquery.FindOne(s.document, s.path)
	rS := htmlquery.InnerText(node)
	if rS == "" {
		return ErrNotFound
	}
	s.returnValue = rS
	return nil
}

func (s *Scraper) RunQuery() (string, error) {
	err := s.GetDocument()
	if err != nil {
		return "", err
	}
	err = s.FindElements()
	if err != nil {
		return "", err
	}
	return s.returnValue, nil
}
