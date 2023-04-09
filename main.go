package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"git.sr.ht/~vicentereyes/ultrastar/transmission"
	"github.com/mmcdole/gofeed"
)

var target string
var alreadyAdded map[string]bool

func main() {
	log.SetFlags(0)
	url := os.Getenv("ULTRASTAR_RSS")
	if url == "" {
		log.Fatal("ULTRASTAR_RSS not set")
	}
	target = os.Getenv("ULTRASTAR_TARGET")

	var err error
	alreadyAdded, err = transmission.List()
	if err != nil {
		log.Fatalf("Error: getting transmission state: %v", err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Printf("Error: downloading RSS feed: %v", err)
		fmt.Println("Attempting to read from STDIN...")
		feed, err = fp.Parse(os.Stdin)
		if err != nil {
			log.Fatalf("Error: reading RSS from stdin: %v", err)
		}
	}
	fmt.Println("RSS obtained successfully!")

	for _, song := range feed.Items {
		addTorrent(song)
	}
}

var regex = regexp.MustCompile(`^(.*?) #`)

func addTorrent(song *gofeed.Item) {
	match := regex.FindStringSubmatch(song.Title)
	if match != nil && alreadyAdded[match[1]] {
		log.Printf("skipping %s: already added", match[1])
		return
	}
	for _, enc := range song.Enclosures {
		if enc.Type != "application/x-bittorrent" {
			continue
		}
		if err := transmission.Add(enc.URL, target); err != nil {
			log.Fatalf("Error: adding %s: %v", song.Title, err)
		}
		return
	}
	log.Printf("Error: adding %s: no torrents available", song.Title)
}
