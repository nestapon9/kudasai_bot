// envío de foto + título + link, y plan B sin imagen

package telegram

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func SendPosts(c *Client, logger *slog.Logger, posts []feed.Post, feedStore string) error {

	var logs []any
	numPosts := len(posts)
	i := 0

	for ; i <= numPosts-1; i++ {
		if c.Ctx.Err() != nil {
			state.NewLogError(logger, time.Now(), fmt.Sprintf("Sent %d/%d. Failed to send %d", i, numPosts, numPosts-i), errors.New("Bot closed"))
			break
		}
		p := posts[i]
		p.TotalPost = numPosts
		p.PostNum = i
		msg := fmt.Sprintf("<a href='%v'>%v</a>", p.URL, p.Title)
		if feedStore == "/animenetwork" {
			err := c.SendMessage(logger, p, msg, feedStore)
			if err != nil {
				state.NewLogError(logger, time.Now(), fmt.Sprintf("Sent %d/%d. Failed to send (title) {%s}", i, numPosts, p.Title), err)
				break
			} else {
				logs = append(logs, SuccessfulSent(p, feedStore))
			}
		} else {
			if err := c.SendPhoto(logger, p, feedStore); err != nil {
				logger.Warn(fmt.Sprintf("Failed to send (title) {%s}. Trying just text message", p.Title), "Warning_date", time.Now(), "Error", err)
				if err = c.SendMessage(logger, p, msg, feedStore); err != nil {
					state.NewLogError(logger, time.Now(), fmt.Sprintf("Sent %d/%d. Failed to send (title) {%s}", i, numPosts, p.Title), err)
					break
				} else {
					logs = append(logs, SuccessfulSent(p, feedStore))
				}
			} else {
				logs = append(logs, SuccessfulSent(p, feedStore))
			}
		}
	}

	if feedStore == "/kudasai" {
		state.WriteStateFile(c.StateKudasai, logger, logs)
	} else {
		state.WriteStateFile(c.StateAnn, logger, logs)
	}

	return nil
}

func SuccessfulSent(p feed.Post, feedStore string) struct {
	Feed string            `json:"Feed"`
	Sent map[string]string `json:"Sent"`
} {
	return struct {
		Feed string            `json:"Feed"`
		Sent map[string]string `json:"Sent"`
	}{
		Feed: feedStore,
		Sent: map[string]string{
			"ID":      p.ID,
			"PubDate": p.PubDate.String(),
		},
	}
}
