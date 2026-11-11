package commands

import (
	"log/slog"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func FetchAll(c *telegram.Client, logger *slog.Logger, posts []feed.Post, feedStore string, feedState []filter.StateLog) error {
	filter.OrderPostsByDate(posts)
	if feedStore == "/animenetwork" {
		if len(posts) > 20 {
			posts = posts[len(posts)-21:]
		}
	}
	toSent := filter.FilterBySent(logger, posts, feedState)
	err := telegram.SendPosts(c, logger, toSent, feedStore)
	if err != nil {
		return err
	}
	return nil
}
