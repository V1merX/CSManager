package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type replyKeyboardValue string

const (
	ReplyServers = replyKeyboardValue("Сервера")
)

func (b *Bot) StartCmd(upd tgbotapi.Update) {
	message := `
	Добро пожаловать в бота!
Выберите нужную категорию:
	`

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, message)
	reply.ParseMode = tgbotapi.ModeMarkdown

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сервера", marshallCb(callbackEntity{cbType: 1})),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление привилегиями", marshallCb(callbackEntity{cbType: 5})),
		),
	)

	reply.ReplyMarkup = keyboard
	reply.DisableWebPagePreview = true

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send start message", zap.Error(err))
	}
}