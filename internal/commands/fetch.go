package commands

import (
	"log/slog"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func Fetch(c *telegram.Client, logger *slog.Logger, category, feedStore string, posts []feed.Post, feedState []filter.StateLog) error {
	filteredPosts := filter.FilterByCategory(posts, category)

	lastestPosts := filter.FilterBySent(logger, filteredPosts, feedState)
	if len(lastestPosts) == 0 {
		logger.Info("Nothing to fetch")
		c.SendMessage(logger, feed.Post{}, "Nothing to fetch", "")
		return nil
	}

	filter.OrderPostsByDate(lastestPosts)

	if err := telegram.SendPosts(c, logger, lastestPosts, feedStore); err != nil {
		return err
	}

	return nil

}
