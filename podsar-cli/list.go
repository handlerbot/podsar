package main

import (
	"fmt"
	"strings"

	"github.com/handlerbot/podsar/lib"
)

func listCmd(db *lib.PodsarDb) error {
	feeds, err := db.GetAllFeeds(false)
	if err != nil {
		return err
	}
	for _, f := range feeds {
		c, err := db.GetEpisodeCount(f.Id)
		if err != nil {
			return err
		}
		s := fmt.Sprintf("* \"%s\" [%s]: %d saved episodes", f.FeedName, f.OurName, c)
		opts := f.SummarizeOptions()
		if len(opts) > 0 {
			s += fmt.Sprintf(" (%s)", strings.Join(opts, ", "))
		}
		fmt.Println(s)
	}
	return nil
}
