package lib

import (
	sql "database/sql"
	"errors"
	"fmt"
	"net/url"
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

func DirAndFilename(itemUrl string, title string, prefix string, f *Feed) (dir, fn string, err error) {
	if f.RenameEpisodesToTitle {
		fn = title + ".mp3"
	} else {
		u, err := url.Parse(itemUrl)
		if err != nil {
			return "", "", errors.New(fmt.Sprintf("unable to parse URL \"%s\": %s", itemUrl, err.Error()))
		}
		fn = filepath.Base(u.Path)
		if fn == "." || fn == string(filepath.Separator) {
			return "", "", errors.New(fmt.Sprintf("unable to extract filename from URL path \"%s\"", fn))
		}
	}
	dir = filepath.Join(prefix, f.DirName)
	fn = filepath.Join(dir, fn)
	return
}
