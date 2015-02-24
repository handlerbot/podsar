package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/handlerbot/podsar/lib"
)

type guidCache struct {
	db    *lib.PodsarDb
	cache map[int]map[string]int
	sync.Mutex
}

func NewGuidCache(db *lib.PodsarDb, ch *chan os.Signal) *guidCache {
	c := guidCache{}
	c.db = db
	c.cache = make(map[int]map[string]int)
	if ch != nil {
		go func() {
			for true {
				select {
				case s := <-*ch:
					c.Lock()
					log.Printf("Received %s signal, flushing seen episodes cache\n", s)
					c.cache = make(map[int]map[string]int)
					c.Unlock()
				}
			}
		}()
	}
	return &c
}

func (c *guidCache) getMapForFeed(id int) (map[string]int, error) {
	m, ok := c.cache[id]
	if !ok {
		if eps, err := c.db.GetAllEpisodes(id); err != nil {
			return nil, errors.New(fmt.Sprintf("loading seen episodes for feed id %d: %s", id, err))
		} else {
			m := make(map[string]int)
			for _, e := range eps {
				m[e.Guid] = 1
			}
			return m, nil
		}
	}
	return m, nil
}

func (c *guidCache) Seen(id int, guid string) (seen bool, err error) {
	c.Lock()
	defer c.Unlock()

	m, err := c.getMapForFeed(id)
	if err != nil {
		return false, err
	}
	_, seen = m[guid]
	return seen, nil
}

func (c *guidCache) MarkSeen(id int, guid string) (err error) {
	c.Lock()
	defer c.Unlock()

	m, err := c.getMapForFeed(id)
	if err != nil {
		return err
	}
	m[guid] = 1
	return
}
