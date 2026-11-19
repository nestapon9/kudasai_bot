// lectura del XML y conversión de fechas a UTC
package feed

import (
	"log/slog"
	"time"
)

func PubDateParser(logger *slog.Logger, pubDate string, title string, timeFormat string) time.Time {
	t, err := time.Parse(timeFormat, pubDate)
	if err != nil {
		logger.Error("Error parsing date", "Post", title, "Error", err)
		return time.Time{}
	}
	return t
}
