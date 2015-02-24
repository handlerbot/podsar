package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dustin/go-humanize"

	"github.com/handlerbot/podsar/lib"
)

var (
	limitStr = flag.String("bwlimit", "0", "limit podcast downloads to this many bytes per second, e.g. 500000, 256KB, 512KiB, 1MB, 3MiB (case insensitive)")
	dbfn     = flag.String("db", "podsar.db", "filename of our sqlite3 database")
	dir      = flag.String("dir", "", "base of directory tree to store downloaded podcasts in")
	temp     = flag.String("temp", "", "temporary directory for in-flight downloads; if not set, defaults to \"(dir)/.podsar-tmp\". MUST BE ON THE SAME FILESYSTEM AS --storage-base!")
)

func main() {
	flag.Parse()

	if *dir == "" {
		bail("base podcast directory must be set via --dir\n")
	}

	if *temp == "" {
		*temp = fmt.Sprintf("%s/.podsar-tmp", *dir)
		log.Println("Using temporary directory", *temp)
	}

	if err := os.MkdirAll(*dir, 0755); err != nil {
		bail("error creating podcast directory \"%s\": %s\n", *dir, err)
	}

	if err := os.MkdirAll(*temp, 0755); err != nil {
		bail("error creating temporary directory \"%s\": %s\n", *temp, err)
	}

	db, err := lib.NewPodsarDb(*dbfn)
	if err != nil {
		bail(fmt.Sprintf("error opening database \"%s\": %s\n", *dbfn, err))
	}
	defer db.Close()

	var limit int64
	if x, err := humanize.ParseBytes(*limitStr); err != nil {
		bail(err.Error())
	} else {
		limit = int64(x)
	}
	if limit > 0 {
		log.Printf("Limiting download bandwidth to %s bytes per second (%s)\n", humanize.Comma(limit), humanize.IBytes(uint64(limit)))
	}

	var wg sync.WaitGroup
	ch := make(chan *retrieveRequest, 100)
	alarm, hup := make(chan os.Signal, 1), make(chan os.Signal, 1)
	quit := make(chan struct{})

	cache := NewGuidCache(db, &hup)

	go pollFeeds(db, ch, cache, alarm, quit, &wg)
	go retrieve(db, ch, cache, *dir, *temp, limit, quit, &wg)

	signal.Notify(alarm, syscall.SIGALRM)
	signal.Notify(hup, syscall.SIGHUP)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM)
	log.Printf("Received %s signal, beginning shutdown...\n", <-stop)

	close(quit)
	wg.Wait()
	log.Println("Shutdown complete")
}

func bail(s string, m ...interface{}) {
	fmt.Printf("FATAL: "+s, m...)
	os.Exit(1)
}
