package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *Bot) StartCallback(upd tgbotapi.Update, entity callbackEntity) {
	message := `
	Добро пожаловать в бота!
Выберите нужную категорию:
	`

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Сервера", marshallCb(callbackEntity{cbType: 1})),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Управление привилегиями", marshallCb(callbackEntity{cbType: 5})),
			),
		),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send servers message", zap.Error(err))
	}
}