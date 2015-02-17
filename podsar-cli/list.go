package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/handlerbot/podsar/lib"
)

func printFeed(f lib.Feed) string {
	s := fmt.Sprintf("%s [%s]", f.FeedName, f.OurName)
	attrs := make([]string, 0)
	if !f.Active {
		attrs = append(attrs, "paused")
	}
	if len(f.DirName) > 0 {
		attrs = append(attrs, fmt.Sprintf("directory \"%s\"", f.DirName))
	}
	if f.RenameEpisodesToTitle {
		attrs = append(attrs, "rename to title")
	}
	if len(attrs) > 0 {
		s += fmt.Sprintf(" (%s)", strings.Join(attrs, ", "))
	}
	return s
}

func listCmd(db *lib.PodsarDb) {
	if feeds, err := db.GetFeeds(true); err != nil {
		log.Fatal(err)
	} else {
		for _, feed := range feeds {
			fmt.Println("*", printFeed(*feed))
		}
	}
}