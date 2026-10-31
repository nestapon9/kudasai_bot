package bot

import (
	"log/slog"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func CommandRun(logger *slog.Logger, c telegram.Client, commandMsg string) error {
	timeNewsMarker, err := filter.GetTimeMarker(logger)
	if err != nil {
		return err
	}

	posts, err := feed.Feed(c.Clt, logger)
	if err != nil {
		return err
	}

	if commandMsg == "fetch all" {
		filter.OrderPostsByDate(posts)
		err = telegram.SendPosts(c, logger, posts)
		if err != nil {
			return err
		}
		return nil
	}

	filteredPosts := filter.FilterByCategory(posts)

	filter.OrderPostsByDate(filteredPosts)

	lastestPosts := filter.FilterByDate(logger, timeNewsMarker, filteredPosts)
	if len(lastestPosts) == 0 {
		logger.Info("Nothing to fetch")
		c.SendMessage(logger, feed.Post{}, "Nothing to fetch")
		return nil
	}

	err = telegram.SendPosts(c, logger, lastestPosts)
	if err != nil {
		return err
	}

	return nil
}
