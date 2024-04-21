package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type commandEntity struct {
	key commandKey
	desc string
	action func(upd tgbotapi.Update)
}

type commandKey string

const (
	StartCmdKey = commandKey("start")
)

func (b *Bot) initCommands() {
	commands := []commandEntity{
		{
			key: StartCmdKey,
			desc: "Запустить бота",
			action: b.StartCmd,
		},
	}

	tgCommands := make([]tgbotapi.BotCommand, 0, len(commands))
	for _, cmd := range commands {
		b.Commands[cmd.key] = cmd
		tgCommands = append(tgCommands, tgbotapi.BotCommand{
			Command:     string(cmd.key),
			Description: cmd.desc,
		})
	}

	tgbotapi.NewSetMyCommands(tgCommands...)
}