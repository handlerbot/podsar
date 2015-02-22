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

func (f *Feed) SummarizeOptions() (a []string) {
	if !f.Active {
		a = append(a, "paused")
	}
	if len(f.DirName) > 0 {
		a = append(a, fmt.Sprintf("directory \"%s\"", f.DirName))
	}
	if f.RenameEpisodesToTitle {
		a = append(a, "rename to title")
	}
	return
}

type ScannerKey struct {
	FeedId int
	Uri    string
}
