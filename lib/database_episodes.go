package lib

import (
	"time"
)

func (p *PodsarDb) PutEpisode(feedId int, title, guid string, timestamp int64) error {
	_, err := p.db.Exec("INSERT INTO episodes(feed_id, title, guid, pub_timestamp) VALUES (?, ?, ?, ?);", feedId, title, guid, timestamp)
	return err
}

func (p *PodsarDb) PutEpisodes(episodes []*Episode) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO episodes(feed_id, title, guid, pub_timestamp) VALUES (?, ?, ?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, e := range episodes {
		_, err = stmt.Exec(e.FeedId, e.Title, e.Guid, e.PubDate.Unix())
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (p *PodsarDb) DeleteAllEpisodes(f *Feed) error {
	_, err := p.db.Exec("DELETE FROM episodes WHERE feed_id = ?;", f.Id)
	return err
}

func (p *PodsarDb) GetAllEpisodes(feedId int) ([]*Episode, error) {
	rows, err := p.db.Query("SELECT title, guid, pub_timestamp FROM episodes WHERE feed_id = ?;", feedId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	episodes := []*Episode{}

	for rows.Next() {
		var title, guid string
		var pubTimestamp int64
		err := rows.Scan(&title, &guid, &pubTimestamp)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, &Episode{feedId, title, guid, time.Unix(pubTimestamp, 0)})
	}
	return episodes, nil
}

func (p *PodsarDb) GetEpisodeCount(feedId int) (int, error) {
	row := p.db.QueryRow("SELECT count(*) FROM episodes WHERE feed_id = ?;", feedId)
	var c int
	if err := row.Scan(&c); err != nil {
		return -1, err
	}
	return c, nil
}
