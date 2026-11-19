package commands

import (
	"log/slog"

	"github.com/arturogood17/kudasai_bot/internal/filter"
	"github.com/arturogood17/kudasai_bot/internal/models"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func FetchAll(c *telegram.Client, logger *slog.Logger, posts []models.Post, feedStore string, feedState []models.StateLog) error {
	filter.OrderPostsByDate(posts)
	toSent := filter.FilterBySent(logger, posts, feedState)
	if len(toSent) == 0 {
		err := c.SendMessage(logger, models.Post{}, "Nothing to fetch", feedStore)
		if err != nil {
			return err
		}
	}
	var finalAllList []models.Post
	for _, s := range toSent {
		if s.Section != "Review" {
			finalAllList = append(finalAllList, s)
		}
	}
	finalAllList = AnimeNetworkLimit(feedStore, finalAllList)
	err := telegram.SendPosts(c, logger, finalAllList, feedStore)
	if err != nil {
		return err
	}
	return nil
}

func AnimeNetworkLimit(feedStore string, finalList []models.Post) []models.Post {
	if feedStore == "/animenetwork" && len(finalList) > 20 {
		return finalList[len(finalList)/2:]
	}
	return finalList
}
