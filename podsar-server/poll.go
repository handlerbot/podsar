package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"

	"github.com/handlerbot/podsar/lib"
)

var (
	feedFetchTimeout = flag.Int("feed-fetch-timeout", 30, "timeout (in seconds) for fetching a podcast's feed")
	pollInterval     = flag.Int("poll-interval", 1800, "delay (in seconds) before checking if we can refresh our feeds")
)

type episodeHandler struct {
	db        *lib.PodsarDb
	cache     *SeenEpisodesCache
	retriever chan *retrieveRequest
	feedId    int
}

func (m *episodeHandler) ProcessItems(scanner *rss.Feed, rssChannel *rss.Channel, rssEntries []*rss.Item) {
	fmt.Println("Processing for feed", m.feedId)
	for _, entry := range rssEntries {
		found, err := m.cache.Seen(m.feedId, *entry.Guid)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !found {
			m.retriever <- &retrieveRequest{m.feedId, entry}
		}
	}
}

func pollFeeds(db *lib.PodsarDb, retrieverCh chan *retrieveRequest, cache *SeenEpisodesCache, trigger chan os.Signal, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	pauseMaybeDone := func(d time.Duration, triggerOK bool) bool {
		wakeup := time.After(d)
		for {
			select {
			case <-wakeup:
				return false
			case s := <-trigger:
				if triggerOK {
					fmt.Printf("Received %s signal, triggering immediate poll\n", s)
					return false
				}
			case <-quit:
				return true
			}
		}
	}

	betweenPolls := time.Duration(time.Duration(*pollInterval) * time.Second)
	midPoll := time.Duration(time.Duration(2) * time.Second)
	feedScanners := make(map[lib.ScannerKey]*rss.Feed)

	for {
		if feeds, err := db.GetAllFeeds(true); err != nil {
			fmt.Println("Error getting list of feeds from database:", err)
		} else {
			for _, feed := range feeds {
				key := lib.ScannerKey{feed.Id, feed.Uri}
				if _, ok := feedScanners[key]; !ok {
					feedScanners[key] = rss.NewWithHandlers(*feedFetchTimeout, true, nil, &episodeHandler{db, cache, retrieverCh, feed.Id})
				}
				scanner := feedScanners[key]
				if !scanner.CanUpdate() {
					// TODO: debug log here?
				} else {
					if err := scanner.Fetch(feed.Uri, nil); err != nil {
						fmt.Printf("Error fetching feed %s (%s): %s\n", feed.OurName, feed.Uri, err)
					}
				}
				if pauseMaybeDone(midPoll, false) {
					return
				}
			}
		}
		fmt.Println("Poller sleeping")
		if pauseMaybeDone(betweenPolls, true) {
			return
		}
	}
}
