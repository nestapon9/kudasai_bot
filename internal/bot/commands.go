package bot

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/commands"
	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

type CMD struct {
	arguments []string
	feedStore string
}

func CommandRun(logger *slog.Logger, c *telegram.Client, commandMsg string) error {
	start := time.Now()
	cmdMsg := CommandProcessing(commandMsg)
	feedState, err := filter.ReadState(c, cmdMsg.feedStore)
	if err != nil {
		return fmt.Errorf("%v - %v", commandMsg, err.Error())
	}

	switch cmdMsg.feedStore {
	case "/kudasai":
		posts, err := feed.Feed[feed.RSS](c.Clt, cmdMsg.feedStore, logger)
		if err != nil {
			return err
		}
		if err := FetchingOptions(c, logger, posts, cmdMsg.feedStore, cmdMsg.arguments, feedState); err != nil { //4. Selección
			return err
		}
		logger.Info("Request received", "Date", start, "Store", cmdMsg.arguments, "Command", commandMsg)
	case "/animenetwork":
		posts, err := feed.Feed[feed.AtomFeed](c.Clt, cmdMsg.feedStore, logger)
		if err != nil {
			return err
		}
		if err := FetchingOptions(c, logger, posts, cmdMsg.feedStore, cmdMsg.arguments, feedState); err != nil {
			return err
		}
		logger.Info("Request received", "Date", start, "Store", cmdMsg.feedStore, "Command", commandMsg)
	}
	return nil
}

func CommandProcessing(command string) CMD {
	splited := strings.Split(strings.TrimSpace(command), " ")
	return CMD{
		feedStore: strings.ToLower(splited[0]),
		arguments: splited[1:],
	}
}

func FetchingOptions(c *telegram.Client, logger *slog.Logger, posts []feed.Post, storeFeed string, arguments []string, feedState []filter.StateLog) error {
	switch len(arguments) {
	case 1:
		if arguments[0] != "all" && arguments[0] != "reset" {
			return errors.New("Command not recognized")
		}
		if arguments[0] == "reset" {
			if err := commands.Reset(c, logger, storeFeed); err != nil {
				return err
			}
		}
		if arguments[0] == "all" {
			if err := commands.FetchAll(c, logger, posts, storeFeed, feedState); err != nil {
				return err
			}
		}

	case 2:
		if arguments[0] != "fetch" && !slices.Contains(config.AcceptedCommands, arguments[1]) {
			return errors.New("Command not recognized")
		}
		if err := commands.Fetch(c, logger, arguments[1], storeFeed, posts, feedState); err != nil {
			return err
		}
	case 3:
		bigThree := strings.Join(arguments[1:], " ")
		if arguments[0] != "fetch" || bigThree != "big three" {
			return errors.New("Command not recognized")
		}
		if err := commands.Fetch(c, logger, bigThree, storeFeed, posts, feedState); err != nil {
			return err
		}
	}
	return nil
}
