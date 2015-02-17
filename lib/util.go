package lib

import (
	sql "database/sql"
	"path/filepath"
)

type scannableRow interface {
	Scan(...interface{}) error
}

func makeFeedFromRow(row scannableRow) (*Feed, error) {
	var id int
	var ourName, feedName, uri string
	var dirName sql.NullString
	var active, rename bool
	err := row.Scan(&id, &ourName, &feedName, &uri, &active, &dirName, &rename)
	if err != nil {
		return nil, err
	}
	return &Feed{id, ourName, feedName, uri, active, dirName.String, rename}, nil
}

func AssembleDest(srcUrl string, title string, dirPrefix string, feed *Feed) (string, string) {
	urlFilename := filepath.Base(srcUrl) // eeek

	var feedDir, destFilename string

	if feed.DirName != "" {
		feedDir = feed.DirName
	} else {
		feedDir = feed.OurName
	}

	if feed.RenameEpisodesToTitle {
		destFilename = title + ".mp3"
	} else {
		destFilename = urlFilename
	}

	var p, f string

	if len(dirPrefix) > 0 {
		p = filepath.Join(dirPrefix, feedDir)
	} else {
		p = feedDir
	}

	f = filepath.Join(p, destFilename)

	return p, f
}
