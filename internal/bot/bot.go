package bot

import (
	"CSManager/internal/config"
	"CSManager/internal/storage/mysql"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Bot struct {
	API       *tgbotapi.BotAPI
	Config    *config.Config
	Logger    *zap.Logger
	Bucket 	  map[int64]actionEntity
	Commands  map[commandKey]commandEntity
	Callbacks map[callbackType]callbackFn
	Storage  *mysql.Storage
}

func NewBot(botApi *tgbotapi.BotAPI, cfg *config.Config, log *zap.Logger) *Bot {
	b := &Bot{
		API:       botApi,
		Config:    cfg,
		Logger:    log,
		Commands:  make(map[commandKey]commandEntity),
		Callbacks: make(map[callbackType]callbackFn),
		Bucket: make(map[int64]actionEntity),
		Storage: nil,
	}

	b.initCommands()
	b.initCallbacks()

	return b
}

func (b *Bot) Run() error {
	b.API.Debug = b.Config.Logger.BotDebug

	b.Logger.Sugar().Infof("Authorized on account %s", b.API.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	defer func() {
		if r := recover(); r != nil {
			b.Logger.Sugar().Errorf("recovered: %s", r)
		}
	}()

	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if b.Config.Logger.BotDebug {
				b.Logger.Sugar().Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
			}

			if !isAdmin(update.Message.From.ID, b.Config.Admins) {
				continue
			}

			if update.Message.IsCommand() {
				key := update.Message.Command()
				if cmd, ok := b.Commands[commandKey(key)]; ok {
					cmd.action(update)
				} else {
					b.Logger.Sugar().Error("command handler not found", zap.String("cmd", key))
				}
			} 

			if b.Bucket[update.Message.From.ID] != (actionEntity{}) {
				b.actionForMessage(update)
			}			
		}

		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			entity := unmarshallCb(data)

			if !isAdmin(update.CallbackQuery.From.ID, b.Config.Admins) {
				continue
			}

			b.Callbacks[entity.cbType](update, entity)
		}
	}

	return nil
}

func isAdmin(userID int64, admins map[int64]config.AdminFlags) bool {
	_, ok := admins[userID]
	return ok
}

func backButton(parentType callbackType, parentIds []string) []tgbotapi.InlineKeyboardButton {
	data := callbackEntity{
		parentType: parentType,
		cbType:     Back,
	}
	if len(parentIds) > 0 {
		data.id = parentIds[len(parentIds)-1]
		data.parentIds = parentIds[0 : len(parentIds)-1]
	}
	button := tgbotapi.NewInlineKeyboardButtonData("← Назад", marshallCb(data))

	return []tgbotapi.InlineKeyboardButton{button}
}
