package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	rss "github.com/jteeuwen/go-pkg-rss"

	"github.com/handlerbot/podsar/lib"
)

type feedResponse struct {
	channel *rss.Channel
	items   []*rss.Item
}

func (m *feedResponse) ProcessItems(feed *rss.Feed, channel *rss.Channel, items []*rss.Item) {
	m.channel = channel
	m.items = items
}

func subscribeCmd(db *lib.PodsarDb) (err error) {
	_, err = db.GetFeedByName(*ourName)
	switch err {
	case lib.ErrNoMatchingFeed:
		break
	case nil:
		return errors.New(fmt.Sprintf("short name \"%s\" already in use, please choose another", *ourName))
	default:
		return errors.New("checking short name: " + err.Error())
	}

	if *dirName == "" {
		*dirName = *ourName
	}

	resp := new(feedResponse)
	if err = rss.NewWithHandlers(15, false, nil, resp).Fetch((*(*uri)).String(), nil); err != nil {
		return errors.New("fetch: " + err.Error())
	}

	if *limit < 0 {
		*limit = len(resp.items)
	}

	f := &lib.Feed{0, *ourName, resp.channel.Title, (*(*uri)).String(), false, *dirName, *rename}

	fmt.Printf("Podcast: \"%s\"\nShort name: \"%s\"\nDescription: \"%s\"\n", resp.channel.Title, f.OurName, resp.channel.Description)

	for _, i := range resp.items {
		if e, ok := findAudio(i); ok {
			_, fn, err := lib.DirAndFilename(e.Url, i.Title, "", f)
			if err != nil {
				return errors.New("sample download filename generation: " + err.Error())
			}
			fmt.Printf("Sample download filename: \"%s\"\n", fn)
			break
		}
	}

	var ignore []*rss.Item
	if *limit > 0 {
		if ignore, err = selectAndPrintEpisodes(resp.items, f); err != nil {
			return errors.New("selecting and printing downloadable episodes: " + err.Error())
		}
	} else {
		ignore = resp.items
	}
	fmt.Printf("\nWill mark %d entries as already seen:\n", len(ignore))
	printEpisodes(ignore)

	if *dryrun {
		fmt.Println("Dry run mode: exiting without updating database")
		return nil
	}

	fmt.Printf("\nIf this looks good, type y or yes and hit RETURN to proceed> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	switch strings.ToLower(scanner.Text()) {
	case "y", "yes":
		break
	default:
		fmt.Println("Confirmation not found, exiting")
		return nil
	}

	f.Id, err = db.PutFeed(f)
	if err != nil {
		return errors.New("saving feed: " + err.Error())
	}

	eps := make([]*lib.Episode, 0)
	for _, i := range ignore {
		t, _ := i.ParsedPubDate() // If we can't parse the publication date, default zero-value for time.Time is fine
		eps = append(eps, &lib.Episode{f.Id, i.Title, *i.Guid, t})
	}

	if err = db.PutEpisodes(eps); err != nil {
		return errors.New("saving episodes: " + err.Error())
	}

	if err = db.SetFeedActive(f, true); err != nil {
		return errors.New("unpausing feed: " + err.Error())
	}

	fmt.Println("Subscribed to podcast")
	return nil
}

func printEpisodes(items []*rss.Item) {
	lines := make([][2]string, 0)
	for _, i := range items {
		lines = append(lines, [2]string{"\"" + i.Title + "\"", "(" + pubDateAsString(i) + ")"})
	}
	prettyPrint(lines)
}

func selectAndPrintEpisodes(items []*rss.Item, f *lib.Feed) (ignore []*rss.Item, err error) {
	i, lines := 0, make([][2]string, 0)
	for c := 0; c < *limit && i < len(items); i++ {
		if e, ok := findAudio(items[i]); ok {
			_, fp, err := lib.DirAndFilename(e.Url, items[i].Title, "", f)
			if err != nil {
				return nil, errors.New("computing output filename: " + err.Error())
			}
			lines = append(lines, [2]string{"\"" + items[i].Title + "\"", "(" + pubDateAsString(items[i]) + ") -> \"" + fp + "\""})
			c++
		} else {
			ignore = append(ignore, items[i])
		}
	}
	if i < len(items) {
		ignore = append(ignore, items[i:]...)
	}
	if len(lines) > 0 {
		fmt.Printf("\nWill download %d entries:\n", len(lines))
		prettyPrint(lines)
	}
	return
}

func pubDateAsString(item *rss.Item) string {
	if t, err := item.ParsedPubDate(); err == nil {
		return t.Format("2006-01-02")
	}
	return "<unparseable>"
}

func prettyPrint(lines [][2]string) {
	titleMax := 0
	for _, l := range lines {
		thisLen := len(l[0])
		if thisLen > titleMax {
			titleMax = thisLen
		}
	}
	numMax := len(strconv.Itoa(len(lines)))
	for i, l := range lines {
		fmt.Printf("%[1]*d) %-[3]*s  %s\n", numMax, i+1, titleMax, l[0], l[1])
	}
}

func findAudio(item *rss.Item) (*rss.Enclosure, bool) {
	for _, e := range item.Enclosures {
		if e.Type == "audio/mpeg" {
			return e, true
		}
	}
	return nil, false
}
