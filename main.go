package main

import (
	"fmt"
	"log"
	"os"

	"git.sr.ht/~vicentereyes/ultrastar/transmission"
	"github.com/mmcdole/gofeed"
)

var target string

func main() {
	log.SetFlags(0)
	url := os.Getenv("ULTRASTAR_RSS")
	if url == "" {
		log.Fatal("ULTRASTAR_RSS not set")
	}
	target = os.Getenv("ULTRASTAR_TARGET")

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

func addTorrent(song *gofeed.Item) {
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
