package mysql

import (
	"CSManager/internal/config"
	"CSManager/internal/models"
	"fmt"

	"github.com/Acidic9/go-steam/steamid"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New(dataSourceName config.DB) (*Storage, error) {
	const op = "storage.mysql.New"

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", 
		dataSourceName.User,
		dataSourceName.Password,
		dataSourceName.Host,
		dataSourceName.Port,
		dataSourceName.DBName,
	)), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s Storage) GetAdminsList(serverID string) []models.IKSAdmin {
    var admins []models.IKSAdmin
   	if err := s.db.Table("iks_admins").Where("server_id = ?", serverID).Scan(&admins).Error; err != nil {
        fmt.Println(err)
    }
    return admins
}

func (s Storage) AddAdmin(steamID, name, flags, serverID string, immunity, groupID, end int) {
	admin := models.IKSAdmin{
		SteamID:  steamID,
		Name:     name,
		Flags:    flags,
		Immunity: immunity,
		Group:  groupID,
		TimeEnd:      end,
		ServerID: serverID,
	}

	if err := s.db.Create(&admin).Error; err != nil {
		fmt.Println(err)
	}
}

func (s Storage) DeleteAdmin(steamID, serverID string)  {
	if err := s.db.Where("sid = ? AND server_id = ?", steamID, serverID).Delete(&models.IKSAdmin{}).Error; err != nil {
		fmt.Println(err)
	}
}

func (s Storage) GetVipsList(serverID string) []models.VIP {
    var vips []models.VIP
   	if err := s.db.Table("vip_users").Where("sid = ?", serverID).Scan(&vips).Error; err != nil {
        fmt.Println(err)
    }
    return vips
}

func (s Storage) AddVip(name, group, serverID string, steamID steamid.ID32, endInt int) {
	vip := map[string]interface{}{
		"account_id": steamID,
		"name":      name,
		"group":     group,
		"expires":   endInt,
		"sid":  serverID,
	}

	if err := s.db.Table("vip_users").Create(&vip).Error; err != nil {
		fmt.Println(err)
	}
}

func (s Storage) DeleteVIP(steamID steamid.ID32, serverID string)  {
	if err := s.db.Table("vip_users").Where("sid = ? AND account_id = ?", serverID, steamID).Delete(&models.VIP{}).Error; err != nil {
		fmt.Println(err)
	}
}