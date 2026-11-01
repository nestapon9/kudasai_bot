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

	posts, err := feed.Feed(c.Clt, logger)
	if err != nil {
		return err
	}

	switch commandMsg {
	case "fetch":
		if err := commands.Fetch(c, logger, timeNewsMarker, posts); err != nil {
			return err
		}
	case "fetch all":
		if err := commands.FetchAll(c, logger, posts); err != nil {
			return err
		}
	case "reset":
		if err := commands.Reset(); err != nil {
			return err
		}
	}

	return nil
}
