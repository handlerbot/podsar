package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/mxk/go-flowrate/flowrate"

	"github.com/handlerbot/podsar/lib"
)

type retrieveRequest struct {
	feedId int
	entry  *rss.Item
}

func retrieve(db *lib.PodsarDb, ch chan *retrieveRequest, cache *SeenEpisodesCache, finalDir string, tempDir string, bwlimit int64, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-quit:
			return
		case req := <-ch:
			feed, _ := db.GetFeed(req.feedId)
			for _, enclosure := range req.entry.Enclosures {
				if enclosure.Type == "audio/mpeg" {
					tempFile, err := ioutil.TempFile(tempDir, "")
					if err != nil {
						fmt.Printf("Error creating temporary file in %s: %s\n", tempDir, err)
						continue
					}

					destDir, destFilepath := lib.AssembleDest(enclosure.Url, req.entry.Title, finalDir, feed)
					if err := os.MkdirAll(destDir, 0755); err != nil {
						fmt.Printf("Error making destination directory \"%s\": %s\n", destDir, err)
						continue
					}

					fmt.Println(enclosure.Url, tempFile.Name(), destDir, destFilepath)

					resp, err := http.Get(enclosure.Url)
					if err != nil {
						fmt.Printf("Error retrieving podcast URL \"%s\": %s\n", enclosure.Url, err)
						continue
					}
					defer resp.Body.Close()

					_, err = io.Copy(tempFile, flowrate.NewReader(resp.Body, bwlimit))
					if err != nil {
						fmt.Printf("Error saving podcast URL \"%s\" to temporary file \"%s\": %s", enclosure.Url, tempFile.Name(), err)
						continue
					}

					tempFile.Close()

					if err := os.Rename(tempFile.Name(), destFilepath); err != nil {
						fmt.Printf("Error renaming \"%s\" to \"%s\": %s", tempFile.Name(), destFilepath, err)
						continue
					}

					if err := db.SaveEpisode(req.feedId, *req.entry.Guid); err != nil {
						fmt.Printf("Error saving episode \"%s\" for feed \"%s\"\n", *req.entry.Guid, feed.OurName)
						continue
					}

					if err := cache.MarkSeen(req.feedId, *req.entry.Guid); err != nil {
						fmt.Printf("Error marking as read episode \"%s\" for feed \"%s\"\n", *req.entry.Guid, feed.OurName)
						continue
					}
				}
			}
		}
	}
}
