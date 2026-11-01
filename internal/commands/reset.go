package commands

import (
	"os"
	"path"

	"github.com/arturogood17/kudasai_bot/internal/config"
)

func Reset() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fileToreset := path.Join(cwd, config.DATA, config.LOGFILEPATH)
	if err := os.Truncate(fileToreset, 0); err != nil {
		return err
	}
	return nil
}
