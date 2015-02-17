package main

import (
	"fmt"
	"log"

	rss "github.com/jteeuwen/go-pkg-rss"

	"github.com/handlerbot/podsar/lib"
)

type episodeHandler struct {
	channel *rss.Channel
	items   []*rss.Item
}

func (m *episodeHandler) ProcessItems(scanner *rss.Feed, rssChannel *rss.Channel, rssEntries []*rss.Item) {
	m.channel = rssChannel
	m.items = rssEntries
}

func findAudioEnclosure(item *rss.Item) (*rss.Enclosure, bool) {
	for _, e := range item.Enclosures {
		if e.Type == "audio/mpeg" {
			return e, true
		}
	}
	return nil, false
}

func subscribeCmd() {
	handler := episodeHandler{}
	scanner := rss.NewWithHandlers(15, false, nil, &handler)

	if err := scanner.Fetch((*(*uri)).String(), nil); err != nil {
		log.Fatalf("Error fetching feed %s: %s\n", uri, err)
	}

	feed := &lib.Feed{0, *ourName, handler.channel.Title, (*(*uri)).String(), false, *dirName, *rename}
	ignore := make([]*rss.Item, 0)

	fmt.Printf("### Podcast: \"%s\" (%d entries)\n### Description: \"%s\"\n", handler.channel.Title, len(handler.items), handler.channel.Description)

	if *limit < 0 {
		*limit = len(handler.items)
	}

SamplePrintLoop:
	for _, item := range handler.items {
		if e, ok := findAudioEnclosure(item); ok {
			_, fp := lib.AssembleDest(e.Url, item.Title, "", feed)
			fmt.Printf("### Example downloaded filename: \"%s\"\n", fp)
			break SamplePrintLoop
		}
	}

	if *limit > 0 {
		fmt.Printf("\n### Will download the following entries:\n")

		var i int
		var item *rss.Item
		for c := 0; c < *limit && i < len(handler.items); i++ {
			if e, ok := findAudioEnclosure(handler.items[i]); ok {
				_, fp := lib.AssembleDest(e.Url, handler.items[i].Title, "", feed)
				fmt.Printf("%d) \"%s\" => filename \"%s\"\n", c+1, handler.items[i].Title, fp)
				c++
			} else {
				ignore = append(ignore, item)
			}
		}
		if i < len(handler.items) {
			ignore = append(ignore, handler.items[i:]...)
		}
		fmt.Println()
	} else {
		ignore = handler.items
	}

	fmt.Printf("### Will mark %d entries as already seen\n", len(ignore))
}
