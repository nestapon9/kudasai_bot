// escribir y leer el marcador (líneas fáciles de identificar) y el caso de primera ejecución
package state

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

func WriteStateFile(file *os.File, logger *slog.Logger, logs []any) {
	file.Seek(0, io.SeekEnd)
	encoder := json.NewEncoder(file)

	for _, l := range logs {
		p := l.(struct {
			Feed string            `json:"Feed"`
			Sent map[string]string `json:"Sent"`
		})
		if err := encoder.Encode(p); err != nil {
			logger.Error("Error trying to encode log", "File", file.Name(), "Feed", p.Feed, "PostID", p.Sent["ID"])
		}
	}
	file.Seek(0, io.SeekStart)
}
