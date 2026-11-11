// descarta las noticias anteriores o iguales al marcador y ordena de más antigua a más reciente
package filter

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"slices"
	"sort"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

type StateLog struct {
	FeedStore string
	Sent      map[string]string
}

func OrderPostsByDate(posts []feed.Post) {
	sort.Slice(posts, func(i int, j int) bool {
		return posts[i].PubDate.Before(posts[j].PubDate)
	})
}

func FilterBySent(logger *slog.Logger, posts []feed.Post, feedstate []StateLog) []feed.Post {
	var latestPost []feed.Post
	var IDList []string

	for _, i := range feedstate {
		IDList = append(IDList, i.Sent["ID"])
	}

	for _, p := range posts {
		if !slices.Contains(IDList, p.ID) {
			latestPost = append(latestPost, p)
		}
	}

	return latestPost
}

func ReadState(client *telegram.Client, storeFeed string) ([]StateLog, error) {
	var m []StateLog
	if storeFeed != "/kudasai" && storeFeed != "/animenetwork" {
		return nil, errors.New("Store feed not recognized")
	}

	var file *os.File
	switch storeFeed {
	case "/kudasai":
		file = client.StateKudasai
	case "/animenetwork":
		file = client.StateAnn
	}

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var s StateLog
		if err := decoder.Decode(&s); err != nil {
			return m, err
		}
		m = append(m, s)
	}
	return m, nil
}
