package bot

import (
	"log/slog"

	"github.com/arturogood17/kudasai_bot/internal/commands"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func CommandRun(logger *slog.Logger, c telegram.Client, commandMsg string) error {
	timeNewsMarker, err := filter.GetTimeMarker(logger)
	if err != nil {
		return err
	}

	posts, err := feed.Feed[feed.RSS](c.Clt, logger) //For now, the logic uses RSS. This needs to change depending on the site to get updates from.
	if err != nil {
		return err
	}

	switch commandMsg {
	case "/fetch":
		if err := commands.Fetch(c, logger, timeNewsMarker, posts); err != nil {
			return err
		}
	case "/all":
		if err := commands.FetchAll(c, logger, posts); err != nil {
			return err
		}
	case "/reset":
		if err := commands.Reset(c, logger); err != nil {
			return err
		}
	}

	return nil
}
