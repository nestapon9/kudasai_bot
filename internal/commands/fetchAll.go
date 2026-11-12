package commands

import (
	"log/slog"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func FetchAll(c *telegram.Client, logger *slog.Logger, posts []feed.Post, feedStore string, feedState []filter.StateLog) error {
	filter.OrderPostsByDate(posts)
	toSent := filter.FilterBySent(logger, posts, feedState)
	if len(toSent) == 0 {
		err := c.SendMessage(logger, feed.Post{}, "Nothing to fetch", feedStore)
		if err != nil {
			return err
		}
	}
	var finalAllList []feed.Post
	for _, s := range toSent {
		if s.Section != "Review" {
			finalAllList = append(finalAllList, s)
		}
	}
	if feedStore == "/animenetwork" && len(finalAllList) > 20 {
		finalAllList = finalAllList[len(finalAllList)-21:]
	}
	err := telegram.SendPosts(c, logger, finalAllList, feedStore)
	if err != nil {
		return err
	}
	return nil
}
