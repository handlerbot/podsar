package main

import (
	"log"

	"github.com/handlerbot/podsar/lib"
)

func pauseUnpauseCmd(db *lib.PodsarDb, name string, state bool) {
	f, err := db.GetFeedByName(name)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.PutFeedState(f, state); err != nil {
		log.Fatal(err)
	}
}