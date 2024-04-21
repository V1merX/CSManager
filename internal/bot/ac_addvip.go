package bot

import (
	"CSManager/internal/storage/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"github.com/Acidic9/go-steam/steamid"
)

func (b *Bot) AddVip(upd tgbotapi.Update) {
	dbServerID := b.Bucket[upd.Message.From.ID].server_id
	dbName := "vip"

	dbConfig := b.Config.Databases[b.Config.Servers[dbServerID].Databases[dbName].DB]

	storage, err := mysql.New(dbConfig)
	if err != nil {
		b.Logger.Sugar().Errorf("failed to create MySQL storage: %s", err)
		return
	}

	serverID := b.Config.Servers[dbServerID].Databases[dbName].ServerID

	vipData := upd.Message.Text

	parts := strings.Split(vipData, ",")

	if len(parts) != 4 {
		b.Logger.Sugar().Error("Invalid number of arguments")
		return
	}

	steamID := strings.TrimSpace(parts[0])
	name := strings.TrimSpace(parts[1])
	group := strings.TrimSpace(parts[2])
	end := strings.TrimSpace(parts[3])

	endInt, err := strconv.Atoi(end)
	if err != nil {
		b.Logger.Sugar().Error("Error converting end to int:", zap.Error(err))
		return
	}

	steamIDInt, err := strconv.ParseUint(steamID, 10, 64)
    if err != nil {
		b.Logger.Sugar().Error("Error converting steamID to int:", zap.Error(err))
        return
    }

	accountID := steamid.NewID64(steamIDInt).To3().To32()

	storage.AddVip(name, group, serverID, accountID, endInt)

	reply := tgbotapi.NewMessage(upd.Message.Chat.ID, "✅ VIP-игрок успешно добавлен")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		backButton(5, []string{}),
	)

	reply.ReplyMarkup = keyboard

	b.Bucket[upd.Message.From.ID] = actionEntity{}

	if _, err := b.API.Send(reply); err != nil {
		b.Logger.Sugar().Error("Failed to send adminadd message:", zap.Error(err))
	}
}
