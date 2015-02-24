package main

import (
	"flag"
	"log"
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

type feedPoller struct {
	db        *lib.PodsarDb
	cache     *guidCache
	retriever *audioRetriever
	scanners  map[lib.ScannerKey]*rss.Feed
	reducers  map[int]*feedReducer
	pause     time.Duration
	trigger   chan os.Signal
	quit      chan struct{}
}

func newFeedPoller(db *lib.PodsarDb, c *guidCache, r *audioRetriever, t chan os.Signal) (p *feedPoller) {
	p = new(feedPoller)
	p.db = db
	p.cache = c
	p.retriever = r
	p.trigger = t

	p.scanners = make(map[lib.ScannerKey]*rss.Feed, 0)
	p.reducers = make(map[int]*feedReducer, 0)

	p.pause = time.Duration(time.Duration(*pollInterval) * time.Second)

	return
}

func (p *feedPoller) Poll(quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	shutdown := false

	go func() {
		select {
		case <-quit:
			shutdown = true
		}
	}()

	for {
		if shutdown {
			return
		}

		if feeds, err := p.db.GetAllFeeds(true); err != nil {
			log.Println("Error getting feeds from database:", err)
		} else {
			for _, f := range feeds {
				var ok bool

				var reducer *feedReducer
				if reducer, ok = p.reducers[f.Id]; !ok {
					reducer = newFeedReducer(f.Id, p.cache)
					p.reducers[f.Id] = reducer
				}

				var scanner *rss.Feed
				key := lib.ScannerKey{f.Id, f.Uri}
				if scanner, ok = p.scanners[key]; !ok {
					scanner = rss.NewWithHandlers(*feedFetchTimeout, true, nil, reducer)
					p.scanners[key] = scanner
				}

				if err = scanner.Fetch(f.Uri, nil); err != nil {
					log.Printf("Error fetching feed \"%s\" (\"%s\"): %s\n", f.Uri, f.OurName, err)
				}

				for _, i := range reducer.Items {
					if shutdown {
						return
					}
					if err = p.retriever.Download(f.Id, i); err != nil {
						log.Printf("Error downloading podcast \"%s\" (\"%s\"): %s\n", i.Title, f.OurName, err)
						continue
					}
					log.Printf("Downloaded \"%s\" (\"%s\")\n", i.Title, f.FeedName)
				}

				reducer.Reset()
			}
		}

		wakeup := time.After(p.pause)
		select {
		case <-quit:
			return
		case <-wakeup:
		case s := <-p.trigger:
			log.Printf("Received %s signal, triggering immediate poll\n", s)
		}
	}
}
