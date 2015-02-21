package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/handlerbot/podsar/lib"
)

func listCmd(db *lib.PodsarDb) {
	if all, err := db.GetAllFeeds(false); err != nil {
		log.Fatal(err)
	} else {
		for _, f := range all {
			if c, err := db.GetEpisodeCount(f.Id); err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("*", printFeed(*f, c))
			}
		}
	}
}

func printFeed(f lib.Feed, count int) string {
	s := fmt.Sprintf("[%s] %s: %d known episodes", f.OurName, f.FeedName, count)
	a := make([]string, 0)
	if !f.Active {
		a = append(a, "paused")
	}
	if len(f.DirName) > 0 {
		a = append(a, fmt.Sprintf("directory \"%s\"", f.DirName))
	}
	if f.RenameEpisodesToTitle {
		a = append(a, "rename to title")
	}
	if len(attrs) > 0 {
		s += fmt.Sprintf(" (%s)", strings.Join(a, ", "))
	}
	return s
}
