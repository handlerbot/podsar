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

func markDone(db *lib.PodsarDb, cache *SeenEpisodesCache, f *lib.Feed, req *retrieveRequest) bool {
	t, _ := req.entry.ParsedPubDate()
	if err := db.PutEpisode(req.feedId, req.entry.Title, *req.entry.Guid, t.Unix()); err != nil {
		fmt.Printf("Error saving episode \"%s\" for feed \"%s\"\n", req.entry.Title, f.OurName)
		return false
	}
	if err := cache.MarkSeen(req.feedId, *req.entry.Guid); err != nil {
		fmt.Printf("Error marking as read episode \"%s\" for feed \"%s\"\n", req.entry.Title, f.OurName)
		return false
	}
	return true
}

func retrieve(db *lib.PodsarDb, ch chan *retrieveRequest, cache *SeenEpisodesCache, finalDir string, tempDir string, bwlimit int64, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-quit:
			return
		case req := <-ch:
			feed, _ := db.GetFeedById(req.feedId)
			found := false

			for _, enclosure := range req.entry.Enclosures {
				if enclosure.Type == "audio/mpeg" {
					found = true

					tempFile, err := ioutil.TempFile(tempDir, "")
					if err != nil {
						fmt.Printf("Error creating temporary file in %s: %s\n", tempDir, err)
						break
					}
					defer tempFile.Close()

					destDir, destFilepath := lib.AssembleDest(enclosure.Url, req.entry.Title, finalDir, feed)
					if err := os.MkdirAll(destDir, 0755); err != nil {
						fmt.Printf("Error making destination directory \"%s\": %s\n", destDir, err)
						break
					}

					fmt.Println(enclosure.Url, tempFile.Name(), destDir, destFilepath)

					resp, err := http.Get(enclosure.Url)
					if err != nil {
						fmt.Printf("Error retrieving podcast URL \"%s\": %s\n", enclosure.Url, err)
						break
					}
					defer resp.Body.Close()

					_, err = io.Copy(tempFile, flowrate.NewReader(resp.Body, bwlimit))
					if err != nil {
						fmt.Printf("Error saving podcast URL \"%s\" to temporary file \"%s\": %s", enclosure.Url, tempFile.Name(), err)
						break
					}

					if err := os.Rename(tempFile.Name(), destFilepath); err != nil {
						fmt.Printf("Error renaming \"%s\" to \"%s\": %s", tempFile.Name(), destFilepath, err)
						break
					}

					_ = markDone(db, cache, feed, req)
					break
				}
			}

			if !found {
				_ = markDone(db, cache, feed, req)
			}
		}
	}
}
