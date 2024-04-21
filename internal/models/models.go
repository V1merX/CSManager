package models

type IKSAdmin struct {
	ID int    `gorm:"column:id"`
	SteamID string `gorm:"column:sid"`
	Name string `gorm:"column:name"`
	Flags string `gorm:"column:flags"`
	Immunity int `gorm:"column:immunity"`
	Group int `gorm:"column:group_id"`
	TimeEnd int `gorm:"column:end"`
	ServerID string `gorm:"column:server_id"`
}

type VIP struct {
	AccountID int    `gorm:"column:account_id"`
	Name string `gorm:"column:name"`
	LastVisit int `gorm:"column:lastvisit"`
	ServerID string `gorm:"column:sid"`
	Group string `gorm:"column:group"`
	Expires int `gorm:"column:expires"`
}