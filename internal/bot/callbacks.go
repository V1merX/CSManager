package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type callbackType int

const (
	Start callbackType = iota
	Servers
	Server 
	Players
	Rcon
	PrivilegiesMenu
	AdminMenu
	ChooseAdminServer
	AdminsList
	AddAdmin
	DeleteAdmin
	VIPMenu
	ChooseVipServer
	VipsList
	AddVip
	DeleteVip
	Back
)

type (
	callbackEntity struct {
		cbType     callbackType
		id         string
		parentType callbackType
		parentIds  []string
		server_id	string
		page int
		user_id	string
	}

	callbackFn func(upd tgbotapi.Update, entity callbackEntity)
)

func (c callbackEntity) Clone() callbackEntity {
	return callbackEntity{
		cbType:     c.cbType,
		id:         c.id,
		parentType: c.parentType,
		parentIds:  c.parentIds,
		server_id:       c.server_id,
		page: c.page,
		user_id:         c.user_id,
	}
}

func (b *Bot) initCallbacks() {
	b.Callbacks = map[callbackType]callbackFn{
		Servers:        b.ServersCallback,
		Server:			b.ServerCallback,
		Back:			b.BackCallback,
		Start:			b.StartCallback,
		Players: 		b.PlayersCallback,
		Rcon:			b.RconCallback,
		PrivilegiesMenu:		b.PrivilegiesMenuCallback,
		AdminMenu: b.AdminMenuCallback,
		AdminsList: b.AdminsListCallback,
		ChooseAdminServer: b.ChooseAdminServerCallback,
		AddAdmin: b.AddAdminCallback,
		DeleteAdmin: b.DeleteAdminCallback,
		VIPMenu: b.VIPMenuCallback,
		ChooseVipServer: b.ChooseVipServerCallback,
		VipsList: b.VipsListCallback,
		AddVip: b.AddVipCallback,
		DeleteVip: b.DeleteVipCallback,
	}
}

func marshallCb(data callbackEntity) string {
	return fmt.Sprintf(
		"%d;%s;%d;%s;%s;%d",
		data.cbType,
		data.id,
		data.parentType,
		strings.Join(data.parentIds, "."),
		data.server_id,
		data.page,
	)
}

func unmarshallCb(data string) callbackEntity {
	d := strings.Split(data, ";")

	var cbType int
	if len(d) > 0 {
		cbType, _ = strconv.Atoi(d[0])
	}

	var id string
	if len(d) > 1 {
		id = d[1]
	}

	var pType int
	if len(d) > 2 {
		pType, _ = strconv.Atoi(d[2])
	}

	var parentIds []string
	if len(d) > 3 {
		parentIds = strings.Split(d[3], ".")
	}

	var server_id string
	if len(d) > 4 {
		server_id = d[4]
	}

	var page int
	if len(d) > 5 {
		page, _ = strconv.Atoi(d[5])
	}

	var user_id string
	if len(d) > 6 {
		user_id = d[6]
	}

	return callbackEntity{
		cbType:     callbackType(cbType),
		id:         id,
		parentType: callbackType(pType),
		parentIds:  parentIds,
		server_id: server_id,
		page: page,
		user_id: user_id,
	}
}