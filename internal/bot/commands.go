package bot

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/commands"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/files"
	"github.com/arturogood17/kudasai_bot/internal/models"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

type CMD struct {
	arguments []string
	feedStore string
}

func CommandRun(logger *slog.Logger, c *telegram.Client, commandMsg string) error {
	start := time.Now()
	cmdMsg := CommandProcessing(commandMsg)
	var posts []models.Post
	var err error
	var file *os.File
	switch cmdMsg.feedStore {
	case "/kudasai":
		file = c.StateKudasai
	case "/animenetwork":
		file = c.StateAnn
	}
	feedState, err := files.ReadState(file, cmdMsg.feedStore)
	if err != nil {
		return fmt.Errorf("%v - %v", commandMsg, err.Error())
	}

	switch cmdMsg.feedStore {
	case "/kudasai":
		posts, err = feed.Feed[feed.RSS](c.Clt, cmdMsg.feedStore, logger)
		if err != nil {
			return err
		}
		if err := FetchingOptions(c, logger, posts, cmdMsg.feedStore, cmdMsg.arguments, feedState); err != nil {
			return err
		}
		logger.Info("Request received", "Date", start, "Store", cmdMsg.arguments, "Command", commandMsg)
	case "/animenetwork":
		posts, err = feed.Feed[feed.AtomFeed](c.Clt, cmdMsg.feedStore, logger)
		if err != nil {
			return err
		}
		if err := FetchingOptions(c, logger, posts, cmdMsg.feedStore, cmdMsg.arguments, feedState); err != nil {
			return err
		}
		logger.Info("Request received", "Date", start, "Store", cmdMsg.feedStore, "Command", commandMsg)
	}

	files.StateTriming(file, posts, logger, cmdMsg.feedStore)

	return nil
}

func CommandProcessing(command string) CMD {
	splited := strings.Split(strings.TrimSpace(command), " ")
	return CMD{
		feedStore: strings.ToLower(splited[0]),
		arguments: splited[1:],
	}
}

func FetchingOptions(c *telegram.Client, logger *slog.Logger, posts []models.Post, storeFeed string, arguments []string, feedState []models.StateLog) error {
	switch len(arguments) {
	case 1:
		if strings.ToLower(arguments[0]) == "reset" {
			if err := commands.Reset(c, logger, storeFeed); err != nil {
				return err
			}
		} else if strings.ToLower(arguments[0]) == "all" {
			if err := commands.FetchAll(c, logger, posts, storeFeed, feedState); err != nil {
				return err
			}
		} else {
			if err := commands.Fetch(c, logger, strings.ToLower(arguments[0]), storeFeed, posts, feedState); err != nil {
				return err
			}
		}

	case 2:
		bigThree := strings.ToLower(strings.Join(arguments[0:], " "))
		if bigThree != "big three" {
			return errors.New("Command not recognized")
		}
		if err := commands.Fetch(c, logger, bigThree, storeFeed, posts, feedState); err != nil {
			return err
		}
	}
	return nil
}
