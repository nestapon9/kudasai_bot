// escribir y leer el marcador (líneas fáciles de identificar) y el caso de primera ejecución
package state

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/arturogood17/kudasai_bot/internal/models"
)

func WriteStateFile(file *os.File, logger *slog.Logger, logs []any) {
	encoder := json.NewEncoder(file)

	for _, l := range logs {
		p := l.(models.StateLog)
		if err := encoder.Encode(p); err != nil {
			logger.Error("Error trying to encode log", "File", file.Name(), "Feed", p.Feed, "PostID", p.Sent["ID"])
		}
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		logger.Error("Unable to return to first position of file", "File", file.Name(), "Error", err.Error())
		return
	}
}
