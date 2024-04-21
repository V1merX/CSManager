package bot

import (
	"fmt"
	"html"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorcon/rcon"
	"go.uber.org/zap"
)

func (b *Bot) SendRconMessage(upd tgbotapi.Update) {
	serverID := b.Bucket[upd.Message.From.ID].server_id
	serverConfig, ok := b.Config.Servers[serverID]
	if !ok {
		b.Logger.Sugar().Error("server configuration not found")
		return
	}

	address := fmt.Sprintf("%s:%d", serverConfig.IP, serverConfig.Rcon.Port)
	conn, err := rcon.Dial(address, serverConfig.Rcon.Password)
	if err != nil {
		b.Logger.Sugar().Error("failed to connect to rcon", zap.Error(err))
		return
	}
	defer conn.Close()

	response, err := conn.Execute(upd.Message.Text)
	if err != nil {
		b.Logger.Sugar().Error("failed to send rcon message", zap.Error(err))
		return
	}

	message := fmt.Sprintf("<code>%s</code>", html.EscapeString(response))

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, message)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		backButton(1, []string{}),
	)

	reply.ReplyMarkup = keyboard
	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err != nil {
		b.Logger.Sugar().Error("failed to send rcon message", zap.Error(err))
	}
}
