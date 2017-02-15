package main

import (
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// scrape parses the data in r as html and looks for tags that link to
// other resources.
func scrape(r io.Reader) (map[string]atom.Atom, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	urls := make(map[string]atom.Atom)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.DataAtom {
			case atom.A, atom.Link:
				for _, a := range n.Attr {
					if a.Key == "href" {
						urls[a.Val] = n.DataAtom
						break
					}
				}
			case atom.Audio, atom.Img, atom.Script, atom.Video:
				for _, a := range n.Attr {
					if a.Key == "src" {
						urls[a.Val] = n.DataAtom
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return urls, nil
}
