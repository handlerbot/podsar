package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/mxk/go-flowrate/flowrate"

	"github.com/handlerbot/podsar/lib"
)

type audioRetriever struct {
	db        *lib.PodsarDb
	cache     *guidCache
	dir, temp string
	bwLimit   int64
}

func newAudioRetriever(db *lib.PodsarDb, cache *guidCache, dir, temp string, bwLimit int64) *audioRetriever {
	return &audioRetriever{db, cache, dir, temp, bwLimit}
}

func (r *audioRetriever) Download(id int, item *rss.Item) error {
	f, err := r.db.GetFeedById(id)
	if err != nil {
		return errors.New("looking up feed:" + err.Error())
	}

	e, ok := lib.FindAudio(item)
	if !ok {
		return r.markDone(f, item)
	}

	tempFile, err := ioutil.TempFile(r.temp, "")
	if err != nil {
		return errors.New("creating temporary file: " + err.Error())
	}
	defer tempFile.Close()

	dir, fn, err := lib.DirAndFilename(e.Url, item.Title, r.dir, f)
	if err != nil {
		return errors.New("calculating output filename: " + err.Error())
	}

	if err = os.MkdirAll(dir, 0755); err != nil {
		return errors.New("making final directory \"%s\": " + err.Error())
	}

	resp, err := http.Get(e.Url)
	if err != nil {
		return errors.New("retrieving podcast: " + err.Error())
	}
	defer resp.Body.Close()

	_, err = io.Copy(tempFile, flowrate.NewReader(resp.Body, r.bwLimit))
	if err != nil {
		return errors.New(fmt.Sprintf("copying \"%s\" to temporary file \"%s\": %s", e.Url, tempFile.Name(), err))
	}

	if err = os.Rename(tempFile.Name(), fn); err != nil {
		return errors.New(fmt.Sprintf("renaming \"%s\" -> \"%s\": %s", tempFile.Name(), fn, err))
	}

	return r.markDone(f, item)
}

func (r *audioRetriever) markDone(f *lib.Feed, item *rss.Item) (err error) {
	t, _ := item.ParsedPubDate()
	if err = r.db.PutEpisode(f.Id, item.Title, *item.Guid, t.Unix()); err != nil {
		return errors.New(fmt.Sprintf("saving episode \"%s\" (podcast \"%s\"): %s", item.Title, f.OurName, err.Error()))
	}
	if err = r.cache.MarkSeen(f.Id, *item.Guid); err != nil {
		return errors.New(fmt.Sprintf("marking GUID \"%s\" (podcast \"%s\") as seen: %s", item.Title, f.OurName, err.Error()))
	}
	return
}
