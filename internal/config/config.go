package config

import (
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/state"
	viper "github.com/spf13/viper"
)

const URL = "https://api.telegram.org/bot"
const KUDASAIFEED = "https://somoskudasai.com/noticias/feed/"
const ANNFEED = "https://www.animenewsnetwork.com/all/atom.xml?ann-edition=w"
const ERRORFILEPATH = "ErrorFile.log"
const LOGFILEPATH = "LogInfoFile.log"
const STATEKUDASAI = "state_kudasai.json"
const STATEANN = "state_animenetwork.json"
const INTERNAL = "internal"
const LOGS = "logs"
const DATA = "data"

type Config struct {
	TOKEN_ID string
	CHAT_ID  int
}

func Load(logger *slog.Logger) *Config {
	var config *Config
	viper := viper.New()
	viper.SetConfigFile(".env")

	viper.BindEnv("TOKEN_ID")
	viper.BindEnv("CHAT_ID")

	err := viper.Unmarshal(&config)
	if err != nil {
		state.NewLogError(logger, time.Now(), "Unable to unmarshal", err)
		return nil
	}
	return config
}

func OpenOrCreateLogFile(parentFolder string, filename string) (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(path.Join(cwd, parentFolder), 0750); err != nil {
		return nil, err
	}
	FileCreation := path.Join(cwd, parentFolder, filename)
	File, err := os.OpenFile(FileCreation, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	return File, nil
}
