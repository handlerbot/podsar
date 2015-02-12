package main

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/alecthomas/kingpin.v1"

	"github.com/handlerbot/podsar/lib"
)

var (
	dbfile = kingpin.Flag("dbfile", "filename of our sqlite3 database").Default("podsar.db").String()

	list = kingpin.Command("list", "list all podcasts")
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

func main() {
	cmd := kingpin.Parse()

	db := new(lib.PodsarDb)
	if err := db.Open(*dbfile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch cmd {
	case "list":
		feeds, err := db.GetFeeds(true)
		if err != nil {
			log.Fatal(err)
		}

		for _, feed := range feeds {
			fmt.Println("*", printFeed(*feed))
		}
	}
}
