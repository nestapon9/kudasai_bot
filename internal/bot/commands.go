package bot

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/commands"
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
	cmdMsg := CommandProcessing(commandMsg)                 // 1. Callback
	feedState, err := filter.ReadState(c, cmdMsg.feedStore) // 2. Leer estado
	if err != nil {
		return fmt.Errorf("%v - %v", commandMsg, err.Error())
	}

	switch cmdMsg.feedStore {
	case "/kudasai":
		posts, err := feed.Feed[feed.RSS](c.Clt, "/kudasai", logger) //3. Obtener feed
		if err != nil {
			return err
		}
		if err := FetchingOptions(c, logger, posts, "/kudasai", cmdMsg.arguments, feedState); err != nil { //4. Selección
			return err
		}
		logger.Info("Request received", "Date", start, "Store", cmdMsg.arguments, "Command", commandMsg)
	case "/animenetwork":
		posts, err := feed.Feed[feed.AtomFeed](c.Clt, "/animenetwork", logger)
		if err != nil {
			return err
		}
		if err := FetchingOptions(c, logger, posts, "/animenetwork", cmdMsg.arguments, feedState); err != nil { //4. Selección
			return err
		}
		logger.Info("Request received", "Date", start, "Store", cmdMsg.feedStore, "Command", commandMsg)
	}
	return nil
}

func CommandProcessing(command string) CMD {
	splited := strings.Split(strings.TrimSpace(command), " ")
	return CMD{
		feedStore: splited[0],
		arguments: splited[1:],
	}
}

func FetchingOptions(c *telegram.Client, logger *slog.Logger, posts []feed.Post, storeFeed string, arguments []string, feedState []map[string]any) error {
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
			if err := commands.FetchAll(c, logger, posts, storeFeed); err != nil {
				return err
			}
		}

	case 2:
		if arguments[0] != "fetch" && arguments[1] != "anime" &&
			arguments[1] != "manga" &&
			arguments[1] != "light novels" || arguments[1] != "cultura otaku" {
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
		if err := commands.Fetch(c, logger, arguments[1], storeFeed, posts, feedState); err != nil {
			return err
		}
	}
	return nil
}
