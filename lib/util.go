package lib

import (
	sql "database/sql"
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
