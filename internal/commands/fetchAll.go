package commands

import (
	"log/slog"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func FetchAll(c telegram.Client, logger *slog.Logger, posts []feed.Post) error {
	filter.OrderPostsByDate(posts)
	err := telegram.SendPosts(c, logger, posts)
	if err != nil {
		return err
	}
	return nil
}
