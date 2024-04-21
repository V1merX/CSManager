package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"fmt"
	"github.com/rumblefrog/go-a2s"
	"html"
)

func (b *Bot) ServerCallback(upd tgbotapi.Update, entity callbackEntity) {
    client, err := a2s.NewClient(fmt.Sprintf("%s:%d", b.Config.Servers[entity.server_id].IP, b.Config.Servers[entity.server_id].Port))
    if err != nil {
        b.Logger.Sugar().Errorf("failed to connect to server: %s", err)
        return
    }

    defer client.Close()

    info, err := client.QueryInfo() // QueryInfo, QueryPlayer, QueryRules
    if err != nil {
        b.Logger.Sugar().Errorf("failed to get query-info from server: %s", err)
        return
    }

    var vacStatus rune
    if info.VAC {
        vacStatus = '✓'
    } else {
        vacStatus = '✗'
    }

    message := fmt.Sprintf(
        "<b>Выбранный сервер:</b> <code>%s</code>\n"+
            "<b>IP:</b> <code>%s</code>\n\n"+
            "<b>Текущая карта:</b> <code>%s</code>\n"+
            "<b>Онлайн:</b> <code>%d/%d</code>\n\n"+
            "<b>VAC:</b> <code>%c</code>\n"+
            "<b>Версия:</b> <code>%s</code>\n",
        html.EscapeString(info.Name),
        fmt.Sprintf("%s:%d", b.Config.Servers[entity.server_id].IP, b.Config.Servers[entity.server_id].Port),
        info.Map,
        info.Players,
        info.MaxPlayers,
        vacStatus,
        info.Version,
    )

    players := tgbotapi.NewInlineKeyboardButtonData("Список игроков", marshallCb(callbackEntity{
        cbType:     3,
        server_id:  entity.server_id,
    }))

    rows := make([][]tgbotapi.InlineKeyboardButton, 0, 25)
    if b.Config.Servers[entity.server_id].Rcon.Password != "" && b.Config.Servers[entity.server_id].Rcon.Port != 0 {
    rcon := tgbotapi.NewInlineKeyboardButtonData("Отправить rcon-команду", marshallCb(callbackEntity{
        cbType:     4,
        server_id:  entity.server_id,
    }))
    rows = append(rows, tgbotapi.NewInlineKeyboardRow(players))
        rows = append(rows, tgbotapi.NewInlineKeyboardRow(rcon))
    }
    
    rows = append(rows, backButton(1, []string{}))

    reply := tgbotapi.NewEditMessageTextAndMarkup(
        upd.CallbackQuery.Message.Chat.ID,
        upd.CallbackQuery.Message.MessageID,
        message,
        tgbotapi.NewInlineKeyboardMarkup(rows...),
    )

    reply.ParseMode = "html"

    if _, err := b.API.Send(reply); err != nil {
        b.Logger.Sugar().Error("failed to send servers message", zap.Error(err))
    }
}
