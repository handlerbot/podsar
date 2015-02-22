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
	bwLimitStr = flag.String("bwlimit", "0", "limit podcast downloads to this many bytes per second, e.g. 500000, 256KB, 512KiB, 1MB, 3MiB (case insensitive)")
	dbfn       = flag.String("db", "podsar.db", "filename of our sqlite3 database")
	podcastDir = flag.String("podcast-dir", "", "base of directory tree to store downloaded podcasts in")
	tempDir    = flag.String("temp-dir", "", "temporary directory for in-flight downloads; if not set, defaults to \"(podcast-dir)/.podsar-tmp\". MUST BE ON THE SAME FILESYSTEM AS --storage-base!")
)

func main() {
	flag.Parse()

	if *podcastDir == "" {
		log.Fatal("Base podcast directory must be set via --podcast-dir")
	}

	if *tempDir == "" {
		*tempDir = fmt.Sprintf("%s/.podsar-tmp", *podcastDir)
		fmt.Println("Using temporary directory", *tempDir)
	}

	if err := os.MkdirAll(*podcastDir, 0755); err != nil {
		log.Fatalf("Error creating podcast directory \"%s\": %s\n", *podcastDir, err)
	}

	if err := os.MkdirAll(*tempDir, 0755); err != nil {
		log.Fatalf("Error creating temporary directory \"%s\": %s\n", *tempDir, err)
	}

	var bwLimit int64
	x, err := humanize.ParseBytes(*bwLimitStr)
	if err != nil {
		log.Fatal(err)
	} else {
		bwLimit = int64(x)
	}
	if bwLimit > 0 {
		fmt.Printf("Limiting download bandwidth to %s bytes per second (%s)\n", humanize.Comma(bwLimit), humanize.IBytes(uint64(bwLimit)))
	}

	db, err := lib.NewPodsarDb(*dbfn)
	if err != nil {
		fmt.Printf("Error opening database \"%s\": %s\n", *dbfn, err)
		os.Exit(1)
	}
	defer db.Close()

	seenCache := NewSeenEpisodesCache(db)

	var wg sync.WaitGroup
	alarmTrigger, hupTrigger := make(chan os.Signal, 1), make(chan os.Signal, 1)
	retrieverCh := make(chan *retrieveRequest, 100)
	quit := make(chan struct{})

	go seenCache.Flusher(hupTrigger)
	go pollFeeds(db, retrieverCh, seenCache, alarmTrigger, quit, &wg)
	go retrieve(db, retrieverCh, seenCache, *podcastDir, *tempDir, bwLimit, quit, &wg)

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
