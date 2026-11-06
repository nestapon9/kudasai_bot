// petición HTTP al feed con timeout y comprobación del código 200
package feed

import (
	"encoding/xml"
	"log/slog"
	"net/http"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/state"
)

func Feed[T RSS | AtomFeed](c *http.Client, logger *slog.Logger) ([]Post, error) {
	var feed T
	logger.Info("Feetching started", "Date_fetching_started", time.Now())
	resp, err := c.Get(config.KUDASAIFEED)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Unable to obtain Kudasai feed", err)
		return nil, err
	}
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		state.NewLogError(logger, time.Now(), "Unable parse xml file", err)
		return nil, err
	}

	defer resp.Body.Close()

	return Posts(logger, feed), err

}

func Posts(logger *slog.Logger, posts any) []Post {
	var P []Post
	if rss, ok := posts.(RSS); ok {
		for _, p := range rss.Channel.Item {
			np := Post{
				Title:    p.Title,
				URL:      p.Link,
				PubDate:  PubDateParser(logger, p.PubDate, p.Title, time.RFC1123Z),
				Category: p.Category,
				PhotoURL: p.Content.URL,
			}
			P = append(P, np)
		}
		return P
	}
	atom, _ := posts.(AtomFeed)
	for _, p := range atom.Entry {
		np := Post{
			Title:    p.Title.Text,
			URL:      p.Link.Href,
			PubDate:  p.Published,
			Category: p.Category.Term,
			PhotoURL: "",
		}
		P = append(P, np)
	}
	return P
}
