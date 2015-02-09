package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/handlerbot/podsar/lib"
)

var (
	dbfile = flag.String("dbfile", "podsar.db", "filename of our sqlite3 database")
)

func main() {
	flag.Parse()

	db := new(lib.PodsarDb)
	if err := db.Open(*dbfile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	seenCache := NewSeenEpisodesCache(db)

	var wg sync.WaitGroup
	alarmTrigger, hupTrigger := make(chan os.Signal, 1), make(chan os.Signal, 1)
	retrieverCh := make(chan *retrieveRequest, 100)
	quit := make(chan struct{})

	go seenCache.Flusher(hupTrigger)
	go pollFeeds(db, retrieverCh, seenCache, alarmTrigger, quit, &wg)
	go retrieve(db, retrieverCh, seenCache, quit, &wg)

	signal.Notify(alarmTrigger, syscall.SIGALRM)
	signal.Notify(hupTrigger, syscall.SIGHUP)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	s := <-c

	fmt.Printf("Received %s signal, beginning shutdown... ", s)
	close(quit)
	wg.Wait()
	fmt.Println("done.")
}
