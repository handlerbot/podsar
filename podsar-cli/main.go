package main

import (
	"log"

	"gopkg.in/alecthomas/kingpin.v1"

	"github.com/handlerbot/podsar/lib"
)

var (
	dbfile = kingpin.Flag("dbfile", "filename of our sqlite3 database").Default("podsar.db").String()

	list = kingpin.Command("list", "list all podcasts")

	subscribe = kingpin.Command("sub", "subscribe to a podcast")
	dryRun    = subscribe.Flag("dry-run", "if set, don't write anything to the database, just print").Bool()
	dirName   = subscribe.Flag("dir", "override directory name to download podcast to, inside podcast directory").String()
	rename    = subscribe.Flag("rename-episodes", "override filename from the episode title, rather than whatever the feed gives us").Bool()
	limit     = subscribe.Flag("episode-limit", "download this many episodes from the podcast when subscribing; 0 means none, -1 means all").Default("3").Int()
	ourName   = subscribe.Arg("name", "short name for this podcast").Required().String()
	uri       = subscribe.Arg("uri", "URI for podcast feed").Required().URL()

	unsub = kingpin.Command("unsub", "unsubscribe from a podcast")
	unsubName = unsub.Arg("name", "short name of the podcast to unsubscribe from").Required().String()

	pause = kingpin.Command("pause", "pause downloading of a podcast")
	pauseName = pause.Arg("name", "short name of the podcast to pause").Required().String()

	unpause = kingpin.Command("unpause", "pause downloading of a podcast")
	unpauseName = unpause.Arg("name", "short name of the podcast to unpause").Required().String()
)

func main() {
	cmd := kingpin.Parse()

	db := new(lib.PodsarDb)
	if err := db.Open(*dbfile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch cmd {
	case "list":
		listCmd(db)
	case "sub":
		subscribeCmd(db)
	case "unsub":
		unsubscribeCmd(db)
	case "pause":
		pauseUnpauseCmd(db, *pauseName, false)
	case "unpause":
		pauseUnpauseCmd(db, *unpauseName, true)
	}
}
