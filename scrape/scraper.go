package scrape

type Scraper interface {
	GetDocument(url string) error
	FindElements(path string) error
	RunQuery() (string, error)
}
