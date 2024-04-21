package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *Bot) RconCallback(upd tgbotapi.Update, entity callbackEntity) {
	b.Bucket[upd.CallbackQuery.From.ID] = actionEntity{cbType: entity.cbType, server_id: entity.server_id}

	message := "Введите rcon команду:"

	rows := [][]tgbotapi.InlineKeyboardButton{}
	rows = append(rows, backButton(1, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send rcon message", zap.Error(err))
	}
}