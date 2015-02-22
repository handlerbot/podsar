package lib

import (
	"errors"
)

func fakeBool(x bool) int {
	if x {
		return 1
	}
	return 0
}

func (p *PodsarDb) PutFeed(f *Feed) (id int, err error) {
	if _, err = p.db.Exec("INSERT INTO feeds(our_name, feed_name, uri, active, dir_name, rename_episodes_to_title) "+
		"VALUES (?, ?, ?, ?, ?, ?);", f.OurName, f.FeedName, f.Uri, fakeBool(f.Active), f.DirName,
		f.RenameEpisodesToTitle); err != nil {
		return -1, errors.New("insert: " + err.Error())
	}
	row := p.db.QueryRow("SELECT id FROM feeds WHERE our_name = ?;", f.OurName)
	if err = row.Scan(&id); err != nil {
		return -1, errors.New("select: " + err.Error())
	}
	return id, nil
}

func (p *PodsarDb) DeleteFeed(f *Feed) (err error) {
	if err = p.DeleteAllEpisodes(f); err != nil {
		return errors.New("deleting episodes: " + err.Error())
	}
	if _, err = p.db.Exec("DELETE FROM feeds WHERE id = ?;", f.Id); err != nil {
		return errors.New("deleting feed:" + err.Error())
	}
	return
}

func (p *PodsarDb) SetFeedActive(f *Feed, active bool) (err error) {
	_, err = p.db.Exec("UPDATE feeds SET active = ? WHERE id = ?;", fakeBool(active), f.Id)
	return
}

const (
	getFeedSQL = "SELECT id, our_name, feed_name, uri, active, dir_name, rename_episodes_to_title FROM feeds"
)

func (p *PodsarDb) GetFeedById(id int) (*Feed, error) {
	row := p.db.QueryRow(getFeedSQL+" WHERE id = ?;", id)
	return makeFeedFromRow(row)
}

func (p *PodsarDb) GetFeedByName(name string) (*Feed, error) {
	row := p.db.QueryRow(getFeedSQL+" WHERE our_name = ?;", name)
	return makeFeedFromRow(row)
}

func (p *PodsarDb) GetAllFeeds(activeOnly bool) (feeds []*Feed, err error) {
	rows, err := p.db.Query(getFeedSQL + ";")
	if err != nil {
		return nil, errors.New("select: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		f, err := makeFeedFromRow(rows)
		if err != nil {
			return nil, errors.New("scan: " + err.Error())
		}
		if activeOnly && !f.Active {
			continue
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
}
