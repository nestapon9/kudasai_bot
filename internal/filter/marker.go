// descarta las noticias anteriores o iguales al marcador y ordena de más antigua a más reciente
package filter

import (
	"encoding/json"
	"log/slog"
	"sort"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func OrderPostsByDate(posts []feed.Post) {
	sort.Slice(posts, func(i int, j int) bool {
		return posts[i].PubDate.Before(posts[j].PubDate)
	})
}

func readLogs() ([]map[string]any, error) { //Tiene que cambiar y aceptar argumento por storeFeed
	logFile, err := config.OpenOrCreateLogFile(config.LOGFILEPATH)
	if err != nil {
		return nil, err
	}

	defer logFile.Close()

	var fullFile []map[string]any
	decoder := json.NewDecoder(logFile)

	for decoder.More() {
		var data map[string]any
		err := decoder.Decode(&data)
		if err != nil {
			return nil, err
		}
		if _, ok := data["Last_title_fetch"]; ok {
			fullFile = append(fullFile, data)
		}
	}

	return fullFile, nil
}

func GetTimeMarker(logger *slog.Logger) (time.Time, error) {
	var lastest time.Time

	LogNews, err := readLogs()
	if err != nil {
		state.NewLogError(logger, time.Now(), "Error obtaining the marker", err)
		return lastest, err
	}

	for _, l := range LogNews {
		pTime := feed.PubDateParser(logger, l["LastestFetch"].(string), "News publish marker", time.RFC3339)
		if pTime.After(lastest) {
			lastest = pTime
		}
	}
	return lastest, nil
}

func FilterByDate(logger *slog.Logger, latest time.Time, posts []feed.Post) []feed.Post {
	var lastestPost []feed.Post
	for _, p := range posts {
		if latest.Before(p.PubDate) {
			lastestPost = append(lastestPost, p)
		}
	}
	return lastestPost
}
