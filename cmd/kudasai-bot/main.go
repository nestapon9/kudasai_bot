package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/arturogood17/kudasai_bot/internal/bot"
	"github.com/arturogood17/kudasai_bot/internal/config"
	"github.com/arturogood17/kudasai_bot/internal/state"
	"github.com/arturogood17/kudasai_bot/internal/telegram"
)

func main() {

	LogFile, err := config.OpenOrCreateLogFile(path.Join(config.INTERNAL, config.LOGS), config.LOGFILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR)
	if err != nil {
		fmt.Printf("Couldn't open/create the logging file: %v\n", err)
		os.Exit(1)
	}

	ErrorFile, err := config.OpenOrCreateLogFile(path.Join(config.INTERNAL, config.LOGS), config.ERRORFILEPATH, os.O_APPEND|os.O_CREATE|os.O_RDWR)
	if err != nil {
		fmt.Printf("Couldn't open/create the error logging file: %v\n", err)
		os.Exit(1)
	}

	kudasaiState, err := config.OpenOrCreateLogFile(config.DATA, config.STATEKUDASAI, os.O_CREATE|os.O_RDWR)
	if err != nil {
		fmt.Printf("Couldn't open/create the SomosKudasai state file: %v\n", err)
		os.Exit(1)
	}

	annState, err := config.OpenOrCreateLogFile(config.DATA, config.STATEANN, os.O_CREATE|os.O_RDWR)
	if err != nil {
		fmt.Printf("Couldn't open/create the Anime Network state file: %v\n", err)
		os.Exit(1)
	}

	defer LogFile.Close()
	defer ErrorFile.Close()
	defer kudasaiState.Close()
	defer annState.Close()

	logger := state.NewMultilogger(ErrorFile, LogFile)

	configLoaded := config.Load(logger)
	if configLoaded == nil {
		fmt.Printf("Failed to load secrets: %v\n", err)
		os.Exit(1)
	}

	if configLoaded.CHAT_ID == 0 || configLoaded.TOKEN_ID == "" {
		fmt.Println("Failed to load one of the secrets")
		state.NewLogError(logger, time.Now(), "Failed to load one of the secrets", nil)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client := telegram.New(ctx, configLoaded)

	client.StateAnn = annState
	client.StateKudasai = kudasaiState

	bot.BotLogic(client, logger)

}
