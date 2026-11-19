// envío de foto + título + link, y plan B sin imagen

package telegram

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/models"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func SendPosts(c *Client, logger *slog.Logger, posts []models.Post, feedStore string) error {

	var logs []any
	numPosts := len(posts)
	i := 0

	for ; i <= numPosts-1; i++ {
		if c.Ctx.Err() != nil {
			logger.Error("Sent", fmt.Sprintf("%d/%d", i, numPosts), "Failed to send", slog.IntValue(i-numPosts), "Error", "Bot closed")
			break
		}
		p := posts[i]
		p.TotalPost = numPosts
		p.PostNum = i
		msg := fmt.Sprintf("<a href='%v'>%v</a>", p.URL, p.Title)
		if feedStore == "/animenetwork" {
			err := c.SendMessage(logger, p, msg, feedStore)
			if err != nil {
				logger.Error("Error sending messages", "Sent", fmt.Sprintf("%d/%d", i, numPosts), "Failed to send (title)", p.Title, "Error", err.Error())
				break
			} else {
				logs = append(logs, SuccessfulSent(p, feedStore))
			}
		} else {
			if err := c.SendPhoto(logger, p, feedStore); err != nil {
				logger.Warn(fmt.Sprintf("Failed to send (title) {%s}. Trying just text message", p.Title), "Warning_date", time.Now(), "Error", err)
				if err = c.SendMessage(logger, p, msg, feedStore); err != nil {
					logger.Error("Error sending messages", "Sent", fmt.Sprintf("%d/%d", i, numPosts), "Failed to send (title)", p.Title, "Error", err)
					break
				} else {
					logs = append(logs, SuccessfulSent(p, feedStore))
				}
			} else {
				logs = append(logs, SuccessfulSent(p, feedStore))
			}
		}
	}

	var fileToProcess *os.File

	if feedStore == "/kudasai" {
		fileToProcess = c.StateKudasai
	} else {
		fileToProcess = c.StateAnn
	}
	state.WriteStateFile(fileToProcess, logger, logs)

	return nil
}

func SuccessfulSent(p models.Post, feedStore string) models.StateLog {
	return models.StateLog{
		Feed: feedStore,
		Sent: map[string]string{
			"ID":      p.ID,
			"PubDate": p.PubDate.String(),
		},
	}
}
