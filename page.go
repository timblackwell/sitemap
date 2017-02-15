package main

import "sync"

// Page is a repersentation of a HTML page
type Page struct {
	url       string
	resources []string
	links     []string
	err       error
}

// SafePages safey stores Pages with url key.
// This can be used concurrently as RWMutex is used
// to ensure no reads happen during write.
type SafePages struct {
	v   map[string]*Page
	mux sync.RWMutex
}

// Write page for url
func (c *SafePages) Write(url string, page *Page) {
	c.mux.Lock()
	c.v[url] = page
	c.mux.Unlock()
}

// Read the page stored for given url
func (c *SafePages) Read(url string) (*Page, bool) {
	c.mux.RLock()
	page, exists := c.v[url]
	c.mux.RUnlock()
	return page, exists
}
