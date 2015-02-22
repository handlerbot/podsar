package main

import (
	"fmt"
	"strings"

	"github.com/handlerbot/podsar/lib"
)

func listCmd(db *lib.PodsarDb) error {
	all, err := db.GetAllFeeds(false)
	if err != nil {
		return err
	}
	for _, f := range all {
		c, err := db.GetEpisodeCount(f.Id)
		if err != nil {
			return err
		}
		fmt.Println("*", printFeed(*f, c))
	}
	return nil
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
	if len(a) > 0 {
		s += fmt.Sprintf(" (%s)", strings.Join(a, ", "))
	}

	return s
}
