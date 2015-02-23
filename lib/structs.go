package lib

import (
	"fmt"
	"time"
)

type Episode struct {
	FeedId  int
	Title   string
	Guid    string
	PubDate time.Time
}

type Feed struct {
	Id                    int
	OurName               string
	FeedName              string
	Uri                   string
	Active                bool
	DirName               string
	RenameEpisodesToTitle bool
}

func (f *Feed) SummarizeOptions() (opts []string) {
	if !f.Active {
		opts = append(opts, "paused")
	}
	if f.DirName != f.OurName {
		opts = append(opts, fmt.Sprintf("directory \"%s\"", f.DirName))
	}
	if f.RenameEpisodesToTitle {
		opts = append(opts, "rename to title")
	}
	return
}

type ScannerKey struct {
	FeedId int
	Uri    string
}
