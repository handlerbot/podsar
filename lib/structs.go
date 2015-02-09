package lib

type Episode struct {
	Guid string
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
