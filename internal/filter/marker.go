// descarta las noticias anteriores o iguales al marcador y ordena de más antigua a más reciente
package filter

import (
	"log/slog"
	"slices"
	"sort"

	"github.com/arturogood17/kudasai_bot/internal/models"
)

func OrderPostsByDate(posts []models.Post) {
	sort.Slice(posts, func(i int, j int) bool {
		return posts[i].PubDate.Before(posts[j].PubDate)
	})
}

func FilterBySent(logger *slog.Logger, posts []models.Post, feedstate []models.StateLog) []models.Post {
	var latestPost []models.Post

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
