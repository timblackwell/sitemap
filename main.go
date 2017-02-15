package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("missing seed domain argument")
	}
	t := time.Now()
	pages := SiteMap(os.Args[1], http.Get, scrape, log.Printf)
	for _, page := range pages {
		log.Println()
		if page.err != nil {
			log.Printf("error crawling page: \"%s\": %v\n", page.url, page.err)
			continue
		}
		log.Printf("PAGE %s\n", page.url)
		log.Println("LINKS")
		for _, url := range page.links {
			log.Println(" -- " + url)
		}
		log.Println("RESOURCES")
		for _, url := range page.resources {
			log.Println(" -- " + url)
		}
	}

	log.Printf("Pages %d, elapsed time: %v\n", len(pages), time.Now().Sub(t))
}
