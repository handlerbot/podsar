package main

import (
	"log"

	"github.com/handlerbot/podsar/lib"
)

func unsubscribeCmd(db *lib.PodsarDb) {
	f, err := db.GetFeedByName(*unsubName)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.DeleteFeed(f); err != nil {
		log.Fatal(err)
	}
}