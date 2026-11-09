// petición HTTP al feed con timeout y comprobación del código 200
package feed

import (
	"bytes"
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func Feed(c *http.Client, logger *slog.Logger) ([]Post, error) {
	var rss RSS
	logger.Info("Feetching started", "Date_fetching_started", time.Now())
	resp, err := c.Get(config.FEED)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Unable to obtain Kudasai feed", err)
		return nil, err
	}
	dataBeforeProcessing, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pattern := regexp.MustCompile(`&[[:alnum:]]+;|&`)

	dataAfterProcessing := pattern.ReplaceAll(dataBeforeProcessing, []byte("&amp;"))
	data := bytes.NewReader(dataAfterProcessing)

	decoder := xml.NewDecoder(data)

	if err := decoder.Decode(&rss); err != nil {
		state.NewLogError(logger, time.Now(), "Unable parse xml file", err)
		return nil, err
	}

	defer resp.Body.Close()

	return Posts(logger, rss), err

}

func Posts(logger *slog.Logger, r RSS) []Post {
	var P []Post
	for _, post := range r.Channel.Item {
		np := Post{
			Title:    post.Title,
			URL:      post.Link,
			PubDate:  PubDateParser(logger, post.PubDate, post.Title, time.RFC1123Z),
			Category: post.Category,
			PhotoURL: post.Content.URL,
		}
		P = append(P, np)
	}
	return P
}
