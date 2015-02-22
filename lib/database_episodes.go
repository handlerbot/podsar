package lib

import (
	"errors"
	"time"
)

func (p *PodsarDb) PutEpisode(id int, title, guid string, timestamp int64) (err error) {
	_, err = p.db.Exec("INSERT INTO episodes(feed_id, title, guid, pub_timestamp) VALUES (?, ?, ?, ?);", id, title, guid, timestamp)
	return
}

func (p *PodsarDb) PutEpisodes(eps []*Episode) (err error) {
	tx, err := p.db.Begin()
	if err != nil {
		return errors.New("begin transaction: " + err.Error())
	}
	stmt, err := tx.Prepare("INSERT INTO episodes(feed_id, title, guid, pub_timestamp) VALUES (?, ?, ?, ?);")
	if err != nil {
		return errors.New("prepare: " + err.Error())
	}
	defer stmt.Close()
	for _, e := range eps {
		if _, err = stmt.Exec(e.FeedId, e.Title, e.Guid, e.PubDate.Unix()); err != nil {
			return errors.New("execute: " + err.Error())
		}
	}
	if err = tx.Commit(); err != nil {
		return errors.New("commit: " + err.Error())
	}
	return
}

func (p *PodsarDb) DeleteAllEpisodes(f *Feed) (err error) {
	_, err = p.db.Exec("DELETE FROM episodes WHERE feed_id = ?;", f.Id)
	return
}

func (p *PodsarDb) GetAllEpisodes(id int) (eps []*Episode, err error) {
	rows, err := p.db.Query("SELECT title, guid, pub_timestamp FROM episodes WHERE feed_id = ?;", id)
	if err != nil {
		return nil, errors.New("select: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var title, guid string
		var pubTime int64
		if err = rows.Scan(&title, &guid, &pubTime); err != nil {
			return nil, errors.New("scan: " + err.Error())
		}
		eps = append(eps, &Episode{id, title, guid, time.Unix(pubTime, 0)})
	}
	return
}

func (p *PodsarDb) GetEpisodeCount(id int) (c int, err error) {
	row := p.db.QueryRow("SELECT count(*) FROM episodes WHERE feed_id = ?;", id)
	if err = row.Scan(&c); err != nil {
		return -1, err
	}
	return
}
