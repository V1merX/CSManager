package bot

import (
	"CSManager/internal/storage/mysql"
	"fmt"
	"html"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (b *Bot) PrivilegiesMenuCallback(upd tgbotapi.Update, entity callbackEntity) {

	adminmenu := tgbotapi.NewInlineKeyboardButtonData("Управление администраторами", marshallCb(callbackEntity{
        cbType:     7,
        server_id:  entity.server_id,
    }))
    vipmenu := tgbotapi.NewInlineKeyboardButtonData("Управление VIP-игроками", marshallCb(callbackEntity{
        cbType:    12,
        server_id:  entity.server_id,
    }))

    rows := make([][]tgbotapi.InlineKeyboardButton, 0, 25)
    rows = append(rows, tgbotapi.NewInlineKeyboardRow(adminmenu))
    rows = append(rows, tgbotapi.NewInlineKeyboardRow(vipmenu))
	rows = append(rows, backButton(0, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		"Выберите раздел:",
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send privilegies message", zap.Error(err))
	}
}

func (b *Bot) ChooseAdminServerCallback(upd tgbotapi.Update, entity callbackEntity) {
	buttons := b.serversButtons(6)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(buttons)+3)
	rows = append(rows, buttons...)
	rows = append(rows, backButton(5, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		"Выберите сервер для управления:",
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send adminservers message", zap.Error(err))
	}
}

func (b *Bot) AdminMenuCallback(upd tgbotapi.Update, entity callbackEntity) {
	adminslist := tgbotapi.NewInlineKeyboardButtonData("Список администраторов", marshallCb(callbackEntity{
        cbType:     8,
        server_id:  entity.server_id,
    }))
	addadmin := tgbotapi.NewInlineKeyboardButtonData("Добавить администратора", marshallCb(callbackEntity{
        cbType:     9,
        server_id:  entity.server_id,
    }))
	deleteadmin := tgbotapi.NewInlineKeyboardButtonData("Удалить администратора", marshallCb(callbackEntity{
        cbType:     10,
        server_id:  entity.server_id,
    }))

    rows := make([][]tgbotapi.InlineKeyboardButton, 0, 25)
    rows = append(rows, tgbotapi.NewInlineKeyboardRow(adminslist))
    rows = append(rows, tgbotapi.NewInlineKeyboardRow(addadmin))
    rows = append(rows, tgbotapi.NewInlineKeyboardRow(deleteadmin))
	rows = append(rows, backButton(7, []string{}))

	message := fmt.Sprintf("<b>Выбранный сервер:</b> <code>%s</code>", html.EscapeString(entity.server_id))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send adminmenu message", zap.Error(err))
	}
}

func (b *Bot) AdminsListCallback(upd tgbotapi.Update, entity callbackEntity) {
	const adminsPerPage int = 15

	dbServerID := entity.server_id
	dbName := "iks-admin"

	dbConfig := b.Config.Databases[b.Config.Servers[dbServerID].Databases[dbName].DB]

	storage, err := mysql.New(dbConfig)

	if err != nil {
		b.Logger.Sugar().Errorf("failed to create MySQL storage: %s", err)
		return
	}

	admins := storage.GetAdminsList(b.Config.Servers[entity.server_id].Databases[dbName].ServerID)

	var message string
	message = fmt.Sprintf("<b>Список администраторов -</b> <code>%s</code>\n\n", html.EscapeString(entity.server_id))

	message += "<code>ID | Имя | SteamID | Флаги | Иммунитет</code>\n"

	page := entity.page
	startIndex := page * adminsPerPage
	endIndex := (page + 1) * adminsPerPage
	if endIndex > len(admins) {
		endIndex = len(admins)
	}

	for _, admin := range admins[startIndex:endIndex] {
    	row := fmt.Sprintf("<code>#%d</code> | <code>%s</code> | <code>%s</code> | <code>%s</code> | <code>%d</code>\n", 
        	admin.ID, 
        	html.EscapeString(admin.Name), 
        	html.EscapeString(admin.SteamID), 
        	html.EscapeString(admin.Flags), 
        	admin.Immunity)
    
    	message += row
	}


	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 25)

	// Add navigation buttons
	pagesRow := []tgbotapi.InlineKeyboardButton{}
	if page > 0 {
    	pagesRow = append(pagesRow, tgbotapi.NewInlineKeyboardButtonData("<< Предыдущая",  fmt.Sprintf("8;%s;%d;%s;%s;%d", "", 0, "", entity.server_id, page-1)))
	}

	if endIndex < len(admins) {
		pagesRow = append(pagesRow, tgbotapi.NewInlineKeyboardButtonData("Следующая >>", fmt.Sprintf("8;%s;%d;%s;%s;%d", "", 0, "", entity.server_id, page+1)))
	}
	rows = append(rows, pagesRow)

	rows = append(rows, backButton(7, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send adminslist message", zap.Error(err))
	}
}

func (b *Bot) AddAdminCallback(upd tgbotapi.Update, entity callbackEntity) {
	b.Bucket[upd.CallbackQuery.From.ID] = actionEntity{cbType: entity.cbType, server_id: entity.server_id, user_id: upd.CallbackQuery.From.ID}

	message := "<b>Введите команду:</b>\n\n<i>Пример: </i><code>76561199216836441, Имя, Флаг, Иммунитет, GroupID, Endtime (unixtime)</code>"

	rows := [][]tgbotapi.InlineKeyboardButton{}
	rows = append(rows, backButton(7, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send adminadd message", zap.Error(err))
	}
}

func (b *Bot) DeleteAdminCallback(upd tgbotapi.Update, entity callbackEntity) {
	b.Bucket[upd.CallbackQuery.From.ID] = actionEntity{cbType: entity.cbType, server_id: entity.server_id, user_id: upd.CallbackQuery.From.ID}

	message := "<b>Введите steamid64 пользователя:</b>\n\n<i>Пример: </i><code>76561199216836441</code>"

	rows := [][]tgbotapi.InlineKeyboardButton{}
	rows = append(rows, backButton(7, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send admindelete message", zap.Error(err))
	}
}

func (b *Bot) ChooseVipServerCallback(upd tgbotapi.Update, entity callbackEntity) {
	buttons := b.serversButtons(11)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(buttons)+3)
	rows = append(rows, buttons...)
	rows = append(rows, backButton(5, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		"Выберите сервер для управления:",
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send vipservers message", zap.Error(err))
	}
}

func (b *Bot) VIPMenuCallback(upd tgbotapi.Update, entity callbackEntity) {
	adminslist := tgbotapi.NewInlineKeyboardButtonData("Список VIP-игроков", marshallCb(callbackEntity{
        cbType:     13,
        server_id:  entity.server_id,
    }))
	addadmin := tgbotapi.NewInlineKeyboardButtonData("Добавить VIP-игрока", marshallCb(callbackEntity{
        cbType:     14,
        server_id:  entity.server_id,
    }))
	deleteadmin := tgbotapi.NewInlineKeyboardButtonData("Удалить VIP-игрока", marshallCb(callbackEntity{
        cbType:     15,
        server_id:  entity.server_id,
    }))

    rows := make([][]tgbotapi.InlineKeyboardButton, 0, 25)
    rows = append(rows, tgbotapi.NewInlineKeyboardRow(adminslist))
    rows = append(rows, tgbotapi.NewInlineKeyboardRow(addadmin))
    rows = append(rows, tgbotapi.NewInlineKeyboardRow(deleteadmin))
	rows = append(rows, backButton(7, []string{}))

	message := fmt.Sprintf("<b>Выбранный сервер:</b> <code>%s</code>", html.EscapeString(entity.server_id))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send adminmenu message", zap.Error(err))
	}
}

func (b *Bot) VipsListCallback(upd tgbotapi.Update, entity callbackEntity) {
	const vipsPerPage int = 15

	dbServerID := entity.server_id
	dbName := "vip"

	dbConfig := b.Config.Databases[b.Config.Servers[dbServerID].Databases[dbName].DB]

	storage, err := mysql.New(dbConfig)

	if err != nil {
		b.Logger.Sugar().Errorf("failed to create MySQL storage: %s", err)
		return
	}

	vips := storage.GetVipsList(b.Config.Servers[entity.server_id].Databases[dbName].ServerID)

	var message string
	message = fmt.Sprintf("<b>Список VIP-игроков -</b> <code>%s</code>\n\n", html.EscapeString(entity.server_id))

	message += "<code>AccountID | Имя | Группа | Выдан до</code>\n"

	page := entity.page
	startIndex := page * vipsPerPage
	endIndex := (page + 1) * vipsPerPage
	if endIndex > len(vips) {
		endIndex = len(vips)
	}

	for _, vip := range vips[startIndex:endIndex] {
		row := fmt.Sprintf("<code>%d</code> | <code>%s</code> | <code>%s</code> | <code>%d</code>\n",
			vip.AccountID,
			html.EscapeString(vip.Name),
			html.EscapeString(vip.Group),
			vip.Expires,
		)

		message += row
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	// Add navigation buttons
	pagesRow := []tgbotapi.InlineKeyboardButton{}
	if page > 0 {
    	pagesRow = append(pagesRow, tgbotapi.NewInlineKeyboardButtonData("<< Предыдущая",  fmt.Sprintf("13;%s;%d;%s;%s;%d", "", 0, "", entity.server_id, page-1)))
	}

	if endIndex < len(vips) {
		pagesRow = append(pagesRow, tgbotapi.NewInlineKeyboardButtonData("Следующая >>", fmt.Sprintf("13;%s;%d;%s;%s;%d", "", 0, "", entity.server_id, page+1)))
	}
	rows = append(rows, pagesRow)

	// Add back button
	rows = append(rows, backButton(5, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err != nil {
		b.Logger.Sugar().Error("failed to send vipslist message", zap.Error(err))
	}
}

func (b *Bot) AddVipCallback(upd tgbotapi.Update, entity callbackEntity) {
	b.Bucket[upd.CallbackQuery.From.ID] = actionEntity{cbType: entity.cbType, server_id: entity.server_id, user_id: upd.CallbackQuery.From.ID}

	message := "<b>Введите команду:</b>\n\n<i>Пример: </i><code>76561199216836441, Имя, Группа, Expires (unixtime)</code>"

	rows := [][]tgbotapi.InlineKeyboardButton{}
	rows = append(rows, backButton(5, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send vipadd message", zap.Error(err))
	}
}

func (b *Bot) DeleteVipCallback(upd tgbotapi.Update, entity callbackEntity) {
	b.Bucket[upd.CallbackQuery.From.ID] = actionEntity{cbType: entity.cbType, server_id: entity.server_id, user_id: upd.CallbackQuery.From.ID}

	message := "<b>Введите steamid64 пользователя:</b>\n\n<i>Пример: </i><code>76561199216836441</code>"

	rows := [][]tgbotapi.InlineKeyboardButton{}
	rows = append(rows, backButton(5, []string{}))

	reply := tgbotapi.NewEditMessageTextAndMarkup(
		upd.CallbackQuery.Message.Chat.ID,
		upd.CallbackQuery.Message.MessageID,
		message,
		tgbotapi.NewInlineKeyboardMarkup(rows...),
	)

	reply.ParseMode = "html"

	if _, err := b.API.Send(reply); err!= nil {
		b.Logger.Sugar().Error("failed to send admindelete message", zap.Error(err))
	}
}