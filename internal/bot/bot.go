package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func BotLogic(c telegram.Client, logger *slog.Logger) {
	URL := fmt.Sprintf("%v%v/getUpdates", config.URL, c.Token)

	type requestBody struct {
		Offset  int32 `json:"offset"`
		Timeout int32 `json:"timeout"`
	}

	offset := 0
	retryTime := 1 * time.Second
	countRetries := 0

loop:
	for {
		select {
		case <-c.Ctx.Done():
			state.NewLogClosedInfo(logger, time.Now(), "Bot closed", c.Ctx.Err())
			fmt.Println("\nClosing bot")
			break loop
		default:
			payload := requestBody{
				Offset:  int32(offset),
				Timeout: 30,
			}

			resp, err := telegram.RequestSending(&c, c.Ctx, URL, "POST", payload)
			if errors.Is(err, context.Canceled) {
				state.NewLogClosedInfo(logger, time.Now(), "Bot closed", c.Ctx.Err())
				fmt.Println("\nClosing bot")
				break loop
			} else if err != nil {
				state.NewLogError(logger, time.Now(), "Error trying to send the request for long polling", nil)
				time.Sleep(retryTime)
				retryTime += 1 * time.Second
				if retryTime >= 30*time.Second {
					retryTime = 30 * time.Second
				}
				continue
			}

			retryTime = 1 * time.Second

			if resp.StatusCode >= 400 && resp.StatusCode <= 428 || resp.StatusCode >= 430 && resp.StatusCode <= 499 {
				state.NewLogError(logger, time.Now(), fmt.Sprintf("Client error: %d", resp.StatusCode), nil)
				fmt.Printf("Client error: %d\n", resp.StatusCode)
				fmt.Println("Closing bot")
				return
			} else if resp.StatusCode == 429 || resp.StatusCode >= 500 {
				state.NewLogError(logger, time.Now(), fmt.Sprintf("Server error: %d", resp.StatusCode), nil)
				time.Sleep(retryTime)
				retryTime += 1 * time.Second
				if retryTime >= 30*time.Second {
					retryTime = 30 * time.Second
				}
				continue
			}

			retryTime = 1 * time.Second

			var upd Updates

			if err = json.NewDecoder(resp.Body).Decode(&upd); err != nil {
				state.NewLogError(logger, time.Now(), "Error trying to decode the response from server", err)
				return
			}

			defer resp.Body.Close()

			if len(upd.Result) > 0 {
				mostRecentMsg := upd.Result[len(upd.Result)-1]
				offset = mostRecentMsg.UpdateID + 1
				if !messageAuth(mostRecentMsg.Message.From.ID, c.Chat_id) {
					c.SendMessage(logger, feed.Post{}, "Unavailable")
					state.NewLogError(logger, time.Now(), "Invalid user attempt", nil)
					continue
				}

				CommandMsg := strings.ToLower(mostRecentMsg.Message.Text)

				if CommandMsg != "fetch" && CommandMsg != "fetch all" && CommandMsg != "reset" {
					c.SendMessage(logger, feed.Post{}, "Command not recongnize")
					continue
				}

				if err := CommandRun(logger, c, CommandMsg); err != nil {
					state.NewLogError(logger, time.Now(), "Error serving posts", err)
					offset = mostRecentMsg.UpdateID
					time.Sleep(retryTime)
					retryTime += 1 * time.Second
					if retryTime >= 30*time.Second {
						retryTime = 30 * time.Second
					}
					countRetries += 1
					if countRetries >= 4 {
						offset = mostRecentMsg.UpdateID + 1
						countRetries = 0
					}
					continue
				}
				retryTime = 1 * time.Second
				countRetries = 0
			}
		}

	}
}
