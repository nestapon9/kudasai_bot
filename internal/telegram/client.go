package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"unicode/utf8"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/feed"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

type Client struct {
	Clt          *http.Client
	Ctx          context.Context
	Token        string
	Chat_id      int
	StateKudasai *os.File
	StateAnn     *os.File
}

func New(ctx context.Context, c *config.Config) *Client {
	return &Client{
		Token:   c.TOKEN_ID,
		Chat_id: c.CHAT_ID,
		Clt:     &http.Client{Timeout: 45 * time.Second},
		Ctx:     ctx,
	}
}

func (c *Client) SendMessage(logger *slog.Logger, post feed.Post, message, feedStore string) error {
	type content struct {
		Text      string `json:"text"`
		Chat_ID   int    `json:"chat_id"`
		ParseMode string `json:"parse_mode"`
	}

	URL := fmt.Sprintf("%v%v/sendMessage", config.URL, c.Token)

	payload := content{
		Text:      message,
		Chat_ID:   c.Chat_id,
		ParseMode: "HTML",
	}

	resp, err := RequestSending(c, context.Background(), URL, "POST", payload)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Error trying to send the request", nil)
		return err
	}

	if post != (feed.Post{}) {
		switch feedStore {
		case "/kudasai":
			state.WriteStateFile(c.StateKudasai, logger, map[string]any{
				"feed": feedStore,
				"sent": []map[string]any{
					{
						"id":  post.ID,
						"pub": post.PubDate,
					},
				},
			})
		case "/animenetwork":
			state.WriteStateFile(c.StateAnn, logger, map[string]any{
				"feed": feedStore,
				"sent": []map[string]any{
					{
						"id":  post.ID,
						"pub": post.PubDate,
					},
				},
			})
		}
	}

	logger.Info("Sending posts", "Sent", fmt.Sprintf("%v/%v", post.PostNum+1, post.TotalPost))

	defer resp.Body.Close()

	return nil
}

func (c *Client) SendPhoto(logger *slog.Logger, post feed.Post, feedStore string) error {
	type content struct {
		Photo     string `json:"photo"`
		Caption   string `json:"caption"`
		ParseMode string `json:"parse_mode"`
		Chat_ID   int    `json:"chat_id"`
	}

	URL := fmt.Sprintf("%v%v/sendPhoto", config.URL, c.Token)

	caption := fmt.Sprintf("<a href='%v'>%v</a>", post.URL, post.Title)

	if utf8.RuneCountInString(caption) > 1024 {
		caption = caption[:1024]
	}

	payload := content{
		Photo:     post.PhotoURL,
		Chat_ID:   c.Chat_id,
		Caption:   caption,
		ParseMode: "HTML",
	}

	resp, err := RequestSending(c, context.Background(), URL, "POST", payload)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Error trying to send the request", nil)
		return err
	}

	if post != (feed.Post{}) {
		state.WriteStateFile(c.StateKudasai, logger, map[string]any{
			"feed": feedStore,
			"sent": []map[string]any{
				{
					"id":  post.ID,
					"pub": post.PubDate,
				},
			},
		})
	}

	defer resp.Body.Close()

	return nil

}

func ToMarshalAndBuffer(d any) (*bytes.Buffer, error) {
	p, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(p), nil
}

func RequestSending(c *Client, ctx context.Context, URL string, method string, payload any) (*http.Response, error) {
	data, err := ToMarshalAndBuffer(payload)
	if data == nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, URL, data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	resp, err := c.Clt.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
