package lib

import (
	sql "database/sql"
	"errors"
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
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("no such feed found in database")
	case err != nil:
		return nil, err
	}
	return &Feed{id, ourName, feedName, uri, active, dirName.String, rename}, nil
}

func FinalDirAndFn(url string, title string, prefix string, f *Feed) (finalDir, finalDirAndFn string) {
	var dir, fn string
	if f.RenameEpisodesToTitle {
		fn = title + ".mp3"
	} else {
		fn = filepath.Base(url) // eeek
	}
	if len(prefix) > 0 {
		finalDir = filepath.Join(prefix, f.DirName)
	} else {
		finalDir = dir
	}
	finalDirAndFn = filepath.Join(finalDir, fn)
	return
}
