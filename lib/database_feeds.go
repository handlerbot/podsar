package lib

func fakeBool(x bool) int {
	if x {
		return 1
	}
	return 0
}

func (p *PodsarDb) PutFeed(f *Feed) (int, error) {
	_, err := p.db.Exec("INSERT INTO feeds(our_name, feed_name, uri, active, dir_name, rename_episodes_to_title) "+
		"VALUES (?, ?, ?, ?, ?, ?);", f.OurName, f.FeedName, f.Uri, fakeBool(f.Active), f.DirName,
		f.RenameEpisodesToTitle)
	if err != nil {
		return -1, nil
	}

	row := p.db.QueryRow("SELECT id FROM feeds WHERE our_name = ?;", f.OurName)
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (p *PodsarDb) DeleteFeed(f *Feed) error {
	if err := p.DeleteAllEpisodes(f); err != nil {
		return err
	}
	_, err := p.db.Exec("DELETE FROM feeds WHERE id = ?;", f.Id)
	return err
}

func (p *PodsarDb) PutFeedState(f *Feed, active bool) error {
	_, err := p.db.Exec("UPDATE feeds SET active = ? WHERE id = ?;", fakeBool(active), f.Id)
	return err
}

func (p *PodsarDb) GetFeedById(id int) (*Feed, error) {
	row := p.db.QueryRow("SELECT id, our_name, feed_name, uri, active, dir_name, rename_episodes_to_title FROM feeds WHERE id = ?;", id)
	return makeFeedFromRow(row)
}

func (p *PodsarDb) GetFeedByName(name string) (*Feed, error) {
	row := p.db.QueryRow("SELECT id, our_name, feed_name, uri, active, dir_name, rename_episodes_to_title FROM feeds WHERE our_name = ?;", name)
	return makeFeedFromRow(row)
}

func (p *PodsarDb) GetAllFeeds(activeOnly bool) ([]*Feed, error) {
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

