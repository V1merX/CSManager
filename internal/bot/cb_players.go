package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"fmt"
	"html"
	"github.com/rumblefrog/go-a2s"
)

func (b *Bot) PlayersCallback(upd tgbotapi.Update, entity callbackEntity) {
	client, err := a2s.NewClient(fmt.Sprintf("%s:%d", b.Config.Servers[entity.server_id].IP, b.Config.Servers[entity.server_id].Port))
    if err != nil {
        b.Logger.Sugar().Errorf("failed to connect to server: %s", err)
    }

    defer client.Close()

    info, err := client.QueryPlayer()

    if err != nil {
        b.Logger.Sugar().Errorf("failed to get query-info from server: %s", err)
    }

	var message string
	message = fmt.Sprintf("<b>Список игроков -</b> <code>%s</code>\n\n", html.EscapeString(entity.server_id))
	for i, player := range info.Players {
		message += fmt.Sprintf("<code>#%d -</code> <code>%s</code>\n", i+1, html.EscapeString(player.Name))
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 25)
	rows = append(rows, backButton(1, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send servers message", zap.Error(err))
	}
}