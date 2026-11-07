package bot

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/arturogood17/kudasai_bot/internal/commands"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func CommandRun(logger *slog.Logger, c *telegram.Client, commandMsg string) error {
	timeNewsMarker, err := filter.GetTimeMarker(logger) //A quitar
	if err != nil {
		return err
	}

	feedStore, arguments := CommandProcessing(commandMsg)

	switch feedStore {
	case "/kudasai":
		posts, err := feed.Feed[feed.RSS](c.Clt, logger)
		if err != nil {
			return err
		}
		if err := FetchingOptions(c, logger, posts, "/kudasai", arguments); err != nil {
			return err
		}
	case "/animenetwork":
		posts, err := feed.Feed[feed.AtomFeed](c.Clt, logger)
		if err != nil {
			return err
		}
		if err := FetchingOptions(c, logger, posts, "/animenetwork", arguments); err != nil {
			return err
		}
	}
	return nil
}

func CommandProcessing(command string) (string, []string) {
	splited := strings.Split(strings.TrimSpace(command), " ")
	var processed []string
	for _, s := range splited {
		if s != "" {
			processed = append(processed, s)
		}
	}
	return processed[0], processed[1:]
}

func FetchingOptions(c *telegram.Client, logger *slog.Logger, posts []feed.Post, storeFeed string, arguments []string) error {
	switch len(arguments) {
	case 1:
		if arguments[0] != "fetch" {
			return errors.New("Command not recognized")
		}
		if err := commands.FetchAll(c, logger, posts); err != nil { //Agregar a todas las funciones de fetch la store para poder guardar en el archivo correcto
			return err
		}
	case 2:
		if arguments[0] != "fetch" || arguments[1] != "anime" ||
			arguments[1] != "manga" ||
			arguments[1] != "light novels" {
			return errors.New("Command not recognized")
		}
		if err := commands.FetchByCategory(c, logger, timeNewsMarker, posts); err != nil { //TODO
			return err
		}
	case 3:
		if arguments[0] != "fetch" || strings.Join(arguments[1:], " ") != "big three" {
			return errors.New("Command not recognized")
		}
		if err := commands.Fetch(c, logger, timeNewsMarker, posts); err != nil { //Ver si se puede usar una misma función para 1 categoría o 3. TODO
			return err
		}
	}
}
