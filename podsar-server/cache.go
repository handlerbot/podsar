package main

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/handlerbot/podsar/lib"
)

type SeenEpisodesCache struct {
	db    *lib.PodsarDb
	cache map[int]map[string]int
	sync.Mutex
}

func NewSeenEpisodesCache(db *lib.PodsarDb) *SeenEpisodesCache {
	c := SeenEpisodesCache{}
	c.db = db
	c.cache = make(map[int]map[string]int)

	return &c
}

func (c *SeenEpisodesCache) Flusher(trigger chan os.Signal) {
	for {
		select {
		case s := <-trigger:
			c.Lock()
			fmt.Printf("Received %s signal, flushing seen episodes cache\n", s)
			c.cache = make(map[int]map[string]int)
			c.Unlock()
		}
	}
}

func (c *SeenEpisodesCache) getFeedSeenMap(feedId int) (map[string]int, error) {
	seen, ok := c.cache[feedId]
	if !ok {
		if episodes, err := c.db.GetEpisodes(feedId); err != nil {
			return nil, errors.New(fmt.Sprintf("Error retrieving already seen episodes for feed id %d from database: %s", feedId, err))
		} else {
			seen := make(map[string]int)
			for _, e := range episodes {
				seen[e.Guid] = 1
			}
			return seen, nil
		}
	}
	return seen, nil
}

func (c *SeenEpisodesCache) Seen(feedId int, guid string) (bool, error) {
	c.Lock()
	defer c.Unlock()

	if seen, err := c.getFeedSeenMap(feedId); err != nil {
		return false, err
	} else {
		_, found := seen[guid]
		return found, nil
	}
}

func (c *SeenEpisodesCache) MarkSeen(feedId int, guid string) error {
	c.Lock()
	defer c.Unlock()

	if seen, err := c.getFeedSeenMap(feedId); err != nil {
		return err
	} else {
		seen[guid] = 1
		return nil
	}
}
