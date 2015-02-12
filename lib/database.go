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

func (p *PodsarDb) GetFeed(id int) (*Feed, error) {
	row := p.db.QueryRow("SELECT id, our_name, feed_name, uri, active, dir_name, rename_episodes_to_title FROM feeds WHERE id = ?;", id)
	return makeFeedFromRow(row)
}

func (p *PodsarDb) GetFeeds(activeOnly bool) ([]*Feed, error) {
	rows, err := p.db.Query("SELECT id, our_name, feed_name, uri, active, dir_name, rename_episodes_to_title FROM feeds;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := []*Feed{}

	for rows.Next() {
		feed, err := makeFeedFromRow(rows)
		if err != nil {
			return nil, err
		}
		if activeOnly && !feed.Active {
			continue
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}

func (p *PodsarDb) GetEpisodes(feedId int) ([]*Episode, error) {
	rows, err := p.db.Query("SELECT guid FROM episodes WHERE feed_id = ?;", feedId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	episodes := []*Episode{}

	for rows.Next() {
		var guid string
		err := rows.Scan(&guid)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, &Episode{guid})
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

func (p *PodsarDb) SaveEpisode(feedId int, guid string) error {
	_, err := p.db.Exec("INSERT INTO episodes(feed_id, guid) VALUES (?, ?);", feedId, guid)
	return err
}

func (p *PodsarDb) PauseFeed(feedId int) error {
	_, err := p.db.Exec("UPDATE feeds SET active = 0 WHERE feed_id = ?;", feedId)
	return err
}

func (p *PodsarDb) UnpauseFeed(feedId int) error {
	_, err := p.db.Exec("UPDATE feeds SET active = 1 WHERE feed_id = ?;", feedId)
	return err
}
