package main

import (
	"fmt"
	"sync"

	rss "github.com/jteeuwen/go-pkg-rss"

	"github.com/handlerbot/podsar/lib"
)

type retrieveRequest struct {
	feedId int
	entry  *rss.Item
}

func retrieve(db *lib.PodsarDb, ch chan *retrieveRequest, cache *SeenEpisodesCache, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-quit:
			return
		case req := <-ch:
			feed, _ := db.GetFeed(req.feedId)
			fmt.Println("Saved", feed.OurName, req.entry.Title, *req.entry.Guid)
			db.SaveEpisode(req.feedId, *req.entry.Guid)
			cache.MarkSeen(req.feedId, *req.entry.Guid)
		}
	}
}
