package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v1"

	"github.com/handlerbot/podsar/lib"
)

var (
	dbfn = kingpin.Flag("db", "filename of our sqlite3 database").Default("podsar.db").String()

	list = kingpin.Command("list", "list all podcasts")

	show     = kingpin.Command("show", "show all information about one podcast")
	showName = show.Arg("name", "short name of the podcast to view").Required().String()

	subscribe = kingpin.Command("sub", "subscribe to a podcast")
	dryRun    = subscribe.Flag("dry-run", "if set, don't write anything to the database, just print").Bool()
	dirName   = subscribe.Flag("dir", "directory inside podcast root for this podcast; default is the short name given for the podcast").String()
	rename    = subscribe.Flag("rename-episodes", "create filename from the episode title, rather than using whatever the download URI specifies").Bool()
	limit     = subscribe.Flag("episode-limit", "download this many episodes from the podcast when subscribing; 0 == none").Default("3").Int()
	ourName   = subscribe.Arg("name", "short name for this podcast").Required().String()
	uri       = subscribe.Arg("uri", "URI for podcast feed").Required().URL()

	unsub     = kingpin.Command("unsub", "unsubscribe from a podcast")
	unsubName = unsub.Arg("name", "short name of the podcast to unsubscribe from").Required().String()

	pause     = kingpin.Command("pause", "pause downloading of a podcast")
	pauseName = pause.Arg("name", "short name of the podcast to pause").Required().String()

	resume     = kingpin.Command("resume", "resume downloading of a podcast")
	resumeName = resume.Arg("name", "short name of the podcast to resume").Required().String()
)

func main() {
	cmd := kingpin.Parse()

	db, err := lib.NewPodsarDb(*dbfn)
	if err != nil {
		fmt.Printf("Error opening database \"%s\": %s\n", *dbfn, err)
		os.Exit(1)
	}
	defer db.Close()

	switch cmd {
	case "list":
		err = listCmd(db)
	case "show":
		err = showCmd(db)
	case "sub":
		err = subscribeCmd(db)
	case "unsub":
		err = unsubscribeCmd(db)
	case "pause":
		err = setActiveCmd(db, *pauseName, false)
	case "resume":
		err = setActiveCmd(db, *resumeName, true)
	}

	if err != nil {
		fmt.Printf("Error during \"%s\" command: %s\n", cmd, err)
		os.Exit(1)
	}
}
