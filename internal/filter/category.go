// se queda solo con Anime y Manga
package filter

import (
	"github.com/arturogood17/kudasai_bot/internal/feed"
)

const Manga = "Manga"
const Anime = "Anime"
const Novelas = "Novelas Ligeras"

func FilterByCategory(posts []feed.Post, category string) []feed.Post {
	var filteredPosts []feed.Post
	for _, p := range posts {
		if category != "big three" {
			if p.Category == category && p.Section != "Review" {
				filteredPosts = append(filteredPosts, p)
			}
		} else if (p.Category == Manga || p.Category == Anime || p.Category == Novelas) && p.Section != "Review" {
			filteredPosts = append(filteredPosts, p)
		}
	}
	return filteredPosts
}
