package main

import (
	"fmt"
	"strings"

	"github.com/handlerbot/podsar/lib"
)

func showCmd(db *lib.PodsarDb) (err error) {
	var f *lib.Feed
	if f, err = db.GetFeedByName(*showName); err != nil {
		return
	}

	var c int
	if c, err = db.GetEpisodeCount(f.Id); err != nil {
		return
	}

	fmt.Printf("Podcast Name: \"%s\"\nShort name: \"%s\"\nURI: %s\nSaved episodes: %d\n", f.FeedName, f.OurName, f.Uri, c)

	opts := f.SummarizeOptions()
	if len(opts) > 0 {
		fmt.Printf("Options: %s\n", strings.Join(opts, ", "))
	}

	if c > 0 {
		eps := make([]*lib.Episode, 0)
		if eps, err = db.GetAllEpisodes(f.Id); err != nil {
			return
		}
		fmt.Printf("\nEpisodes:\n")
		for i, e := range eps {
			fmt.Printf("%d) \"%s\"\n", i+1, e.Title)
		}
	}

	return
}
