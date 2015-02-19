package lib

import (
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

type ScannerKey struct {
	FeedId int
	Uri    string
}
