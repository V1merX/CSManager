package main

import (
	"CSManager/internal/bot"
	"CSManager/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	configName string = "local.json"
	configPath string = "../../configs/"
)

func main() {
	cfg, err := config.NewConfig(configName, configPath)
	if err != nil {
		panic(err)
	}

	// server := cfg.Databases["s1"].Password
	// fmt.Println(server)

	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      cfg.Logger.Development,
		Encoding:         cfg.Logger.Encoding,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, _ := zapConfig.Build()
	defer log.Sync()

	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		panic(err)
	}

	bot := bot.NewBot(botApi, cfg, log)
	if err := bot.Run(); err != nil {
		panic(err)
	}
}
