package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	actionEntity struct {
		cbType     callbackType
		server_id	string
		user_id int64
	}
)

func (b *Bot) actionForMessage(upd tgbotapi.Update) {
	switch b.Bucket[upd.Message.From.ID].cbType {
	case 4:
		b.SendRconMessage(upd)
	case 9:
		b.AddAdmin(upd)
	case 10:
		b.DeleteAdmin(upd)
	case 14:
		b.AddVip(upd)
	case 15:
		b.DeleteVip(upd)
	}
}