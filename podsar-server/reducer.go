package main

import (
	"fmt"

	rss "github.com/jteeuwen/go-pkg-rss"
)

type feedReducer struct {
	id    int
	Items []*rss.Item
	Errs  []string
	cache *guidCache
}

func newFeedReducer(id int, cache *guidCache) (fr *feedReducer) {
	fr = new(feedReducer)
	fr.id = id
	fr.cache = cache
	fr.Reset()
	return
}

func (fr *feedReducer) Reset() {
	fr.Items = make([]*rss.Item, 0)
	fr.Errs = make([]string, 0)
}

func (fr *feedReducer) ProcessItems(scanner *rss.Feed, channel *rss.Channel, items []*rss.Item) {
	for _, i := range items {
		found, err := fr.cache.Seen(fr.id, *i.Guid)
		if err != nil {
			fr.Errs = append(fr.Errs, fmt.Sprintf("error checking cache for (\"%s\", %d): %s", *i.Guid, fr.id, err))
			continue
		}
		if !found {
			fr.Items = append(fr.Items, i)
		}
	}
}
