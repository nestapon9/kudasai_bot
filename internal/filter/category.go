// se queda solo con Anime y Manga
package filter

import (
	"github.com/arturogood17/kudasai_bot/internal/models"
)

const Manga = "manga"
const Anime = "anime"
const Novelas = "novelas ligeras"
const Novels = "novels"

func FilterByCategory(posts []models.Post, category string) []models.Post {

	var filteredPosts []models.Post
	for _, p := range posts {
		if category != "big three" {
			if p.Category == category && p.Section != "Review" {
				filteredPosts = append(filteredPosts, p)
			}
		} else if (p.Category == Manga || p.Category == Anime || p.Category == Novelas || p.Category == Novels) && p.Section != "Review" {
			filteredPosts = append(filteredPosts, p)
		}
	}

	return filteredPosts
}
