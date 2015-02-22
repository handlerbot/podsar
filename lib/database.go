package lib

import (
	sql "database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type PodsarDb struct {
	db *sql.DB
}

func NewPodsarDb(fn string) (p *PodsarDb, err error) {
	if _, err = os.Stat(fn); os.IsNotExist(err) {
		return nil, errors.New("database file doesn't exist, please create it (via README.md instructions)")
	}
	p = new(PodsarDb)
	if p.db, err = sql.Open("sqlite3", fn); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *PodsarDb) Close() error {
	return p.db.Close()
}
