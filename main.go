package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Crawler represents a web crawler
type Crawler struct {
	baseURL *url.URL
	visited map[string]bool
}

// NewCrawler creates a new instance of Crawler
func NewCrawler(baseURL string) (*Crawler, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		baseURL: u,
		visited: make(map[string]bool),
	}, nil
}

// Crawl starts crawling from the base URL
func (c *Crawler) Crawl() {
	c.visitPage(c.baseURL.String())
}

func (c *Crawler) visitPage(url string) {
	if c.visited[url] {
		return
	}

	fmt.Println("Visiting:", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error visiting page:", err)
		return
	}
	defer resp.Body.Close()

	c.visited[url] = true

	if resp.StatusCode != http.StatusOK {
		log.Println("Received non-OK status code:", resp.StatusCode)
		return
	}

	links := c.extractLinks(resp.Body)
	for _, link := range links {
		c.visitPage(link)
	}
}

func (c *Crawler) extractLinks(body io.Reader) []string {
	links := make([]string, 0)

	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return links
		}

		token := tokenizer.Token()
		if tokenType == html.StartTagToken && token.Data == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					link := c.resolveURL(attr.Val)
					if link != "" {
						links = append(links, link)
					}
					break
				}
			}
		}
	}
}

func (c *Crawler) resolveURL(href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}

	if u.Scheme == "" {
		u.Scheme = c.baseURL.Scheme
	}
	if u.Host == "" {
		u.Host = c.baseURL.Host
	}

	if c.isSameDomain(u) {
		return u.String()
	}

	return ""
}

func (c *Crawler) isSameDomain(u *url.URL) bool {
	return u.Host == c.baseURL.Host || strings.HasSuffix(u.Host, "."+c.baseURL.Host)
}

func main() {
	crawler, err := NewCrawler("http://example.com")
	if err != nil {
		log.Fatal("Error creating crawler:", err)
	}

	crawler.Crawl()
}
