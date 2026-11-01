package commands

import (
	"log/slog"
	"os"
	"path"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func Reset(c telegram.Client, logger *slog.Logger) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fileToreset := path.Join(cwd, config.DATA, config.LOGFILEPATH)
	if err := os.Truncate(fileToreset, 0); err != nil {
		return err
	}
	if err := c.SendMessage(logger, feed.Post{}, "News reseted"); err != nil {
		return err
	}
	return nil
}
