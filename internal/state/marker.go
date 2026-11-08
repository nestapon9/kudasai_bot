// escribir y leer el marcador (líneas fáciles de identificar) y el caso de primera ejecución
package state

import (
	"encoding/json"
	"log/slog"
	"os"
)

func WriteStateFile(file *os.File, logger *slog.Logger, log map[string]any) {
	data, err := json.Marshal(log)
	if err != nil {
		logger.Error("Error trying to marshal state", "File", file.Name(), "Error", err)
	}
	if _, err := file.Write(data); err != nil {
		logger.Error("Error trying to write state", "File", file.Name(), "Error", err)
	}
}
