package main

import (
	"github.com/handlerbot/podsar/lib"
)

func setActiveCmd(db *lib.PodsarDb, name string, state bool) (err error) {
	var f *lib.Feed
	if f, err = db.GetFeedByName(name); err == nil {
		err = db.SetFeedActive(f, state)
	}
	return
}

func unsubscribeCmd(db *lib.PodsarDb) (err error) {
	var f *lib.Feed
	if f, err = db.GetFeedByName(*unsubName); err == nil {
		err = db.DeleteFeed(f)
	}
	return
}
