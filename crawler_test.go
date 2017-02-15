package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"golang.org/x/net/html/atom"
)

func TestCrawler(t *testing.T) {
	type testFixture struct {
		url     string
		result  map[string]*Page
		httpGet func(string) (*http.Response, error)
		scraper func(io.Reader) (map[string]atom.Atom, error)
	}

	mockErr := fmt.Errorf("unexpected url")
	tests := make(map[string]testFixture)
	tests["happy"] = testFixture{
		url: "https://golang.org/",
		result: map[string]*Page{
			"https://golang.org/": &Page{
				url: "https://golang.org/",
				resources: []string{
					"https://golang.org/assets/stylesheet.css",
					"https://golang.org/assets/script.js",
				},
				links: []string{
					"https://golang.org/",
					"https://golang.org/about",
				},
			},
			"https://golang.org/about": &Page{
				url: "https://golang.org/about",
				err: mockErr,
			},
		},
		httpGet: func(url string) (*http.Response, error) {
			if url != "https://golang.org/" {
				return nil, mockErr
			}
			resp := http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}
			return &resp, nil
		},
		scraper: func(r io.Reader) (map[string]atom.Atom, error) {
			return map[string]atom.Atom{
				"https://golang.org/":      atom.A,
				"https://golang.org/about": atom.A,
				"/assets/stylesheet.css":   atom.Link,
				"/assets/script.js":        atom.Script,
				"https://monzo.me":         atom.A,
			}, nil
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := SiteMap(test.url, test.httpGet, test.scraper, t.Logf)
			if len(result) != len(test.result) {
				t.Fatalf("expected output: '%v', recived: '%v'", test.result, result)
			}
			for url := range test.result {
				if errors.Cause(result[url].err) != test.result[url].err {
					t.Errorf("expected result[%s].err to contain %v, recived: %v",
						url, test.result[url].err, errors.Cause(result[url].err))
				}
				sort.Strings(result[url].links)
				sort.Strings(test.result[url].links)
				if strings.Join(result[url].links, ", ") != strings.Join(test.result[url].links, ", ") {
					t.Errorf("expected links: %s, have %s",
						strings.Join(test.result[url].links, ", "),
						strings.Join(result[url].links, ", "),
					)
				}
				sort.Strings(result[url].resources)
				sort.Strings(test.result[url].resources)
				if strings.Join(result[url].resources, ", ") != strings.Join(test.result[url].resources, ", ") {
					t.Errorf("expected links: %s, have %s",
						strings.Join(test.result[url].resources, ", "),
						strings.Join(result[url].resources, ", "),
					)
				}
			}
		})
	}
}
