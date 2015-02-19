package lib

import (
	sql "database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type PodsarDb struct {
	db *sql.DB
}

func (p *PodsarDb) Open(dbfilename string) error {
	p.Close()
	db, err := sql.Open("sqlite3", dbfilename)
	p.db = db
	return err
}

func (p *PodsarDb) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}
