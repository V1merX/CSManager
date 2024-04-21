package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"fmt"
)

func (b *Bot) serversButtons(cbType int) [][]tgbotapi.InlineKeyboardButton {	
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 25)

	for title, server := range b.Config.Servers {
		ip := server.IP
		port := server.Port

		var data callbackEntity
		data.cbType = callbackType(cbType)
		data.server_id = title

		button := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s (%s:%d)", title, ip, port), marshallCb(data))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	return rows
}

func (b *Bot) ServersCallback(upd tgbotapi.Update, entity callbackEntity) {
	buttons := b.serversButtons(2)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(buttons)+3)
	rows = append(rows, buttons...)
	rows = append(rows, backButton(0, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		"Выберите сервер для настройки:",
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send servers message", zap.Error(err))
	}
}