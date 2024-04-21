
package bot

import (
	"CSManager/internal/storage/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *Bot) DeleteAdmin(upd tgbotapi.Update) {
	dbServerID := b.Bucket[upd.Message.From.ID].server_id
	dbName := "iks-admin"

	dbConfig := b.Config.Databases[b.Config.Servers[dbServerID].Databases[dbName].DB]

	storage, err := mysql.New(dbConfig)
	if err != nil {
		b.Logger.Sugar().Errorf("failed to create MySQL storage: %s", err)
		return
	}

	serverID := b.Config.Servers[dbServerID].Databases[dbName].ServerID

	steamID := upd.Message.Text

	// Add admin to the database
	storage.DeleteAdmin(steamID, serverID)

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, "✅ Администратор успешно удалён")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		backButton(7, []string{}),
	)

	reply.ReplyMarkup = keyboard

	b.Bucket[upd.Message.From.ID] = actionEntity{}

	if _, err := b.API.Send(reply); err != nil {
		b.Logger.Sugar().Error("Failed to send admindelete message:", zap.Error(err))
	}
}

