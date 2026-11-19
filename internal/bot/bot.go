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
	"github.com/arturogood17/kudasai_bot/internal/models"
	"github.com/arturogood17/kudasai_bot/internal/state"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func BotLogic(c *telegram.Client, logger *slog.Logger) {
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

			resp, err := telegram.RequestSending(c, c.Ctx, URL, "POST", payload)
			if errors.Is(err, context.Canceled) {
				state.NewLogClosedInfo(logger, time.Now(), "Bot closed", c.Ctx.Err())
				fmt.Println("\nClosing bot")
				break loop
			} else if err != nil {
				logger.Error("Error trying to send the request for long polling")
				time.Sleep(retryTime)
				retryTime += 1 * time.Second
				if retryTime >= 30*time.Second {
					retryTime = 30 * time.Second
				}
				continue
			}

			retryTime = 1 * time.Second

			if resp.StatusCode >= 400 && resp.StatusCode <= 428 || resp.StatusCode >= 430 && resp.StatusCode <= 499 {
				logger.Error("Client error", "Status code", resp.StatusCode)
				fmt.Printf("Client error: %d\n", resp.StatusCode)
				fmt.Println("Closing bot")
				return
			} else if resp.StatusCode == 429 || resp.StatusCode >= 500 {
				logger.Error("Server error", "Status code", resp.StatusCode)
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
				logger.Error("Error trying to decode the response from server", "Error", err)
				return
			}

			defer resp.Body.Close()

			if len(upd.Result) > 0 {
				mostRecentMsg := upd.Result[len(upd.Result)-1]
				offset = mostRecentMsg.UpdateID + 1
				if !messageAuth(mostRecentMsg.Message.From.ID, c.Chat_id) {
					c.SendMessage(logger, models.Post{}, "Unavailable", "")
					logger.Error("Invalid user attempt", "User", mostRecentMsg.Message.From.ID)
					continue
				}

				CommandMsg := strings.ToLower(mostRecentMsg.Message.Text)

				if !strings.Contains(CommandMsg, "/kudasai") && !strings.Contains(CommandMsg, "/animenetwork") {
					c.SendMessage(logger, models.Post{}, "Feed store not recognized", "")
					continue
				}

				if err := CommandRun(logger, c, CommandMsg); err != nil {
					logger.Error("Error serving posts", "Error", err)
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
