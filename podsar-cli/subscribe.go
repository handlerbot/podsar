package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

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

func subscribeCmd(db *lib.PodsarDb) {
	handler := episodeHandler{}
	scanner := rss.NewWithHandlers(15, false, nil, &handler)

	if err := scanner.Fetch((*(*uri)).String(), nil); err != nil {
		log.Fatalf("Error fetching feed %s: %s\n", uri, err)
	}

	if *limit < 0 {
		*limit = len(handler.items)
	}

	feed := &lib.Feed{0, *ourName, handler.channel.Title, (*(*uri)).String(), false, *dirName, *rename}
	ignore := make([]*rss.Item, 0)

	fmt.Printf("### Podcast: \"%s\" [%s]\n### Description: \"%s\"\n", handler.channel.Title, feed.OurName, handler.channel.Description)

	for _, item := range handler.items {
		if e, ok := findAudioEnclosure(item); ok {
			_, fp := lib.AssembleDest(e.Url, item.Title, "", feed)
			fmt.Printf("### Example downloaded filename: \"%s\"\n", fp)
			break
		}
	}

	fmt.Printf("\n### Found %d entries:\n", len(handler.items))
	maxlen := 0
	lines := make([][2]string, 0)
	for _, item := range handler.items {
		if len(item.Title)+2 > maxlen {
			maxlen = len(item.Title) + 2
		}
		pubDate := "unknown/unparseable publication date"
		if t, err := item.ParsedPubDate(); err == nil {
			pubDate = t.Format("2006-01-02 at 15:04 AM -0700")
		}
		lines = append(lines, [2]string{"\"" + item.Title + "\"", pubDate})
	}
	iWidth := len(strconv.Itoa(len(lines)))
	for i, line := range lines {
		fmt.Printf("%[1]*d) %-[3]*s  (published %s)\n", iWidth, i+1, maxlen, line[0], line[1])
	}

	if *limit > 0 {
		fmt.Printf("\n### Will download the following entries:\n")
		i, maxlen := 0, 0
		lines := make([][2]string, 0)
		for c := 0; c < *limit && i < len(handler.items); i++ {
			if e, ok := findAudioEnclosure(handler.items[i]); ok {
				if len(handler.items[i].Title) > maxlen {
					maxlen = len(handler.items[i].Title) + 2
				}
				_, fp := lib.AssembleDest(e.Url, handler.items[i].Title, "", feed)
				lines = append(lines, [2]string{"\"" + handler.items[i].Title + "\"", fp})
				c++
			} else {
				ignore = append(ignore, handler.items[i])
			}
		}
		if i < len(handler.items) {
			ignore = append(ignore, handler.items[i:]...)
		}
		iWidth := len(strconv.Itoa(len(lines)))
		for i, line := range lines {
			fmt.Printf("%[1]*d) %-[3]*s  => file \"%s\"\n", iWidth, i+1, maxlen, line[0], filepath.Join("<podcast root>", line[1]))
		}
		fmt.Println()
	} else {
		ignore = handler.items
	}

	fmt.Printf("### Will mark %d entries as already seen\n", len(ignore))

	if *dryRun {
		return
	}

	id, err := db.PutFeed(feed)
	if err != nil {
		log.Fatal("Error saving feed:", err)
	}

	episodes := make([]*lib.Episode, 0)
	for _, i := range ignore {
		t, _ := i.ParsedPubDate() // If we can't parse the publication date, default zero-value for time.Time is fine
		episodes = append(episodes, &lib.Episode{id, i.Title, *i.Guid, t})
	}

	if err := db.PutEpisodes(episodes); err != nil {
		log.Fatal("Error saving ignored episodes:", err)
	}

	if err := db.PutFeedState(feed, true); err != nil {
		log.Fatal("Error unpausing feed after creation:", err)
	}
}
