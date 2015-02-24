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
		var eps []*lib.Episode
		if eps, err = db.GetAllEpisodes(f.Id); err != nil {
			return
		}
		fmt.Printf("\nEpisodes:\n")
		g, lines := "", make([][2]string, 0)
		for _, e := range eps {
			if *guid {
				g = fmt.Sprintf(" GUID <%s>", e.Guid)
			}
			lines = append(lines, [2]string{"\"" + e.Title + "\"", g})
		}
		prettyPrint(lines)
	}

	return
}
