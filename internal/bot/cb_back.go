package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (b *Bot) BackCallback(upd tgbotapi.Update, entity callbackEntity) {
	b.Bucket[upd.CallbackQuery.From.ID] = actionEntity{}
	switch entity.parentType {
	case Start:
		b.StartCallback(upd, entity)
	case Servers:
		b.ServersCallback(upd, entity)
	case Server:
		b.ServerCallback(upd, entity)
	case PrivilegiesMenu:
		b.PrivilegiesMenuCallback(upd, entity)
	case AdminMenu:
		b.AdminMenuCallback(upd, entity)
	case ChooseAdminServer:
		b.ChooseAdminServerCallback(upd, entity)
	}
}