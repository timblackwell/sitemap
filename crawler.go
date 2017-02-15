package main

import (
	"io"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html/atom"

	"github.com/pkg/errors"
)

// crawler holds state commen to crawling domain
type crawler struct {
	httpGet func(string) (*http.Response, error)
	scraper func(r io.Reader) (map[string]atom.Atom, error)
	printf  func(format string, v ...interface{})
	domain  string
	pages   *SafePages
}

// HTTPGet makes a HTTP get using url and return the result
type HTTPGet func(url string) (*http.Response, error)

// Scraper gets url-tag map for all links found by parsing HTML in r
type Scraper func(r io.Reader) (map[string]atom.Atom, error)

// Logf to format and log
type Logf func(format string, v ...interface{})

// SiteMap generates collection of pages for domain with resources and links to other
// pages on the domain. An error will be present if there was a problem scraping the page.
func SiteMap(url string, httpGet HTTPGet, scraper Scraper, logger Logf) map[string]*Page {
	pages := make(map[string]*Page)
	crawler := &crawler{
		httpGet: httpGet,
		scraper: scraper,
		printf:  logger,
		domain:  strings.TrimSuffix(url, "/"),
		pages:   &SafePages{v: pages},
	}
	crawler.crawl(url)
	return pages
}

// crawl gets the page repersentation of url and
// then recurses over links found pn page.
// If page already exists for url, the fuction exits.
func (crawler *crawler) crawl(url string) {
	if _, exists := crawler.pages.Read(url); exists {
		return
	}
	crawler.printf("getting page: %s\n", url)
	page := crawler.getPage(url)
	crawler.pages.Write(url, page)
	var wg sync.WaitGroup
	for _, u := range page.links {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			crawler.crawl(url)
		}(u)
	}
	wg.Wait()
}

// getPage will get the repersentation Page for given url filtering out
// links to different domains
func (crawler *crawler) getPage(url string) *Page {
	page := &Page{url: url}
	resp, err := crawler.httpGet(url)
	if err != nil {
		page.err = errors.Wrap(err, "error doing http get")
		return page
	}
	defer resp.Body.Close()
	urls, err := crawler.scraper(resp.Body)
	if err != nil {
		page.err = errors.Wrap(err, "error scraping for links")
		return page
	}

	for u, a := range urls {
		if strings.HasPrefix(u, "/") {
			u = crawler.domain + u
		}
		if a == atom.A {
			if !strings.HasPrefix(u, crawler.domain) {
				continue
			}
			page.links = append(page.links, u)
		} else {
			page.resources = append(page.resources, u)
		}
	}
	return page
}
