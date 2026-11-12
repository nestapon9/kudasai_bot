package commands

import (
	"log/slog"
	"os"
	"path"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func Reset(c *telegram.Client, logger *slog.Logger, feedStore string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	var fileToreset string
	if feedStore == "/kudasai" {
		fileToreset = path.Join(cwd, config.DATA, config.STATEKUDASAI)
	} else {
		fileToreset = path.Join(cwd, config.DATA, config.STATEANN)
	}

	if err := os.Truncate(fileToreset, 0); err != nil {
		return err
	}

	if err := c.SendMessage(logger, feed.Post{}, "News reseted", feedStore); err != nil {
		return err
	}
	return nil
}
