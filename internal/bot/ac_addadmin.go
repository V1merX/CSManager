package bot

import (
	"CSManager/internal/storage/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func (b *Bot) AddAdmin(upd tgbotapi.Update) {
	dbServerID := b.Bucket[upd.Message.From.ID].server_id
	dbName := "iks-admin"

	dbConfig := b.Config.Databases[b.Config.Servers[dbServerID].Databases[dbName].DB]

	storage, err := mysql.New(dbConfig)
	if err != nil {
		b.Logger.Sugar().Errorf("failed to create MySQL storage: %s", err)
		return
	}

	serverID := b.Config.Servers[dbServerID].Databases[dbName].ServerID

	adminData := upd.Message.Text

	parts := strings.Split(adminData, ",")

	if len(parts) != 6 {
		b.Logger.Sugar().Error("Invalid number of arguments")
		return
	}

	steamID := parts[0]
	name := strings.TrimLeft(parts[1], " ")
	flags := strings.TrimLeft(parts[2], " ")
	immunity := strings.TrimSpace(parts[3])
	groupID := strings.TrimSpace(parts[4])
	end := strings.TrimSpace(parts[5])

	immunityInt, err := strconv.Atoi(immunity)
	if err != nil {
		b.Logger.Sugar().Error("Error converting immunity to int:", zap.Error(err))
		return
	}

	groupIDInt, err := strconv.Atoi(groupID)
	if err != nil {
		b.Logger.Sugar().Error("Error converting groupID to int:", zap.Error(err))
		return
	}

	endInt, err := strconv.Atoi(end)
	if err != nil {
		b.Logger.Sugar().Error("Error converting end to int:", zap.Error(err))
		return
	}

	// Add admin to the database
	storage.AddAdmin(steamID, name, flags, serverID, immunityInt, groupIDInt, endInt)

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, "✅ Администратор успешно добавлен")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		backButton(7, []string{}),
	)

	reply.ReplyMarkup = keyboard

	b.Bucket[upd.Message.From.ID] = actionEntity{}

	if _, err := b.API.Send(reply); err != nil {
		b.Logger.Sugar().Error("Failed to send adminadd message:", zap.Error(err))
	}
}

