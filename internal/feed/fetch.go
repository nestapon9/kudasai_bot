// petición HTTP al feed con timeout y comprobación del código 200
package feed

import (
	"bytes"
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/models"
)

func Feed[T RSS | AtomFeed](c *http.Client, storeFeed string, logger *slog.Logger) ([]models.Post, error) {
	var feed T
	logger.Info("Feetching started", "Date_fetching_started", time.Now())
	var urlfeed string
	switch storeFeed {
	case "/kudasai":
		urlfeed = config.KUDASAIFEED
	case "/animenetwork":
		urlfeed = config.ANNFEED
	}
	resp, err := c.Get(urlfeed)
	if err != nil {
		logger.Error("Unable to obtain Kudasai feed", "Error", err)
		return nil, err
	}
	dataBeforeProcessing, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pattern := regexp.MustCompile(`&[[:alnum:]]+;|&`)

	dataAfterProcessing := pattern.ReplaceAll(dataBeforeProcessing, []byte("&amp;"))
	data := bytes.NewReader(dataAfterProcessing)

	if err := xml.NewDecoder(data).Decode(&feed); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	posts, totalNews := Posts(logger, feed)
	logger.Info("Feed fetched", "Total news fetched", totalNews)

	return posts, nil
}

func Posts(logger *slog.Logger, posts any) ([]models.Post, int) {
	var P []models.Post
	if rss, ok := posts.(RSS); ok {
		for _, p := range rss.Channel.Item {
			np := models.Post{
				ID:       p.Link,
				Title:    p.Title,
				URL:      p.Link,
				PubDate:  PubDateParser(logger, p.PubDate, p.Title, time.RFC1123Z),
				Category: strings.ToLower(p.Category),
				PhotoURL: p.Content.URL,
			}
			P = append(P, np)
		}
		return P, len(P)
	}
	atom, _ := posts.(AtomFeed)
	for _, p := range atom.Entry {
		np := models.Post{
			ID:       p.ID.Text,
			Title:    p.Title.Text,
			URL:      p.Link.Href,
			PubDate:  p.Published,
			Category: strings.ToLower(p.Category.Term),
			Section:  p.Cmssection,
			PhotoURL: "",
		}
		P = append(P, np)
	}
	return P, len(P)
}
