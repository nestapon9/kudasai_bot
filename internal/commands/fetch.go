package commands

import (
	"log/slog"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func Fetch(c telegram.Client, logger *slog.Logger, t time.Time, posts []feed.Post) error {
	filteredPosts := filter.FilterByCategory(posts)

	filter.OrderPostsByDate(filteredPosts)

	lastestPosts := filter.FilterByDate(logger, t, filteredPosts)
	if len(lastestPosts) == 0 {
		logger.Info("Nothing to fetch")
		c.SendMessage(logger, feed.Post{}, "Nothing to fetch")
		return nil
	}

	err := telegram.SendPosts(c, logger, posts)
	if err != nil {
		return err
	}

	return nil
}
