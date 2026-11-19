package files

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/arturogood17/kudasai_bot/internal/models"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func ReadState(file *os.File, storeFeed string) ([]models.StateLog, error) {
	var m []models.StateLog
	if storeFeed != "/kudasai" && storeFeed != "/animenetwork" {
		return nil, errors.New("Store feed not recognized")
	}

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var s models.StateLog
		if err := decoder.Decode(&s); err != nil {
			return m, err
		}
		m = append(m, s)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return m, nil
}

func StateTriming(file *os.File, posts []models.Post, logger *slog.Logger, feedStore string) {
	var toBeWritten []any
	latestState, err := ReadState(file, feedStore)
	if err != nil {
		logger.Error("Error reading file state to get latest state", "Feed Store", feedStore, "File", file.Name())
	}

	for _, s := range latestState {
		for _, pf := range posts {
			if s.Sent["ID"] == pf.ID {
				toBeWritten = append(toBeWritten,
					models.StateLog{Feed: feedStore, Sent: map[string]string{"ID": pf.ID, "PubDate": pf.PubDate.String()}})
			}
		}
	}

	if err = os.Truncate(file.Name(), 0); err != nil {
		logger.Error("Error truncating state file", "File", file.Name(), "Error", err.Error())
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		logger.Error("Error returning pointer to first position", "File", file.Name(), "Error", err.Error())
		return
	}

	state.WriteStateFile(file, logger, toBeWritten)
}
