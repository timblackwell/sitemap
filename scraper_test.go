package main

import (
	"bytes"
	"io"
	"testing"

	"golang.org/x/net/html/atom"
)

func TestScraper(t *testing.T) {
	type testFixture struct {
		in  io.Reader
		out map[string]atom.Atom
		err error
	}

	tests := make(map[string]testFixture)
	tests["happy"] = testFixture{
		in: bytes.NewBufferString(
			"<a href=\"https://golang.org/\">golang.org</a>" +
				"<a href=\"https://monzo.me\">monzo.me</a>" +
				"<script async src=\"/assets/script.js\"></script>" +
				"<script async src=\"/assets/script.js\"></script>" +
				"<link rel=\"stylesheet\" type=\"text/css\" href=\"http://googlecode.com/prettify.css\">",
		),
		out: map[string]atom.Atom{
			"https://golang.org/":                atom.A,
			"https://monzo.me":                   atom.A,
			"/assets/script.js":                  atom.Script,
			"http://googlecode.com/prettify.css": atom.Link,
		},
		err: nil,
	}
	tests["empty"] = testFixture{
		in:  bytes.NewBufferString("<p>golang.org</p>"),
		out: nil,
		err: nil,
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out, err := scrape(test.in)
			if err != test.err {
				t.Errorf("expected err: '%v', recived: '%v'", test.err, err)
			}
			if len(out) != len(test.out) {
				t.Fatalf("expected output: '%v', recived: '%v'", test.out, out)
			}
			for i := range test.out {
				if out[i] != test.out[i] {
					t.Errorf("expected output[%s]: '%v', recived: '%v'", i, test.out[i], out[i])
				}
			}
		})
	}
}
