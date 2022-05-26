package Model

import (
	"errors"
	"fmt"
)

type Administrator struct {
	AdminId   string `gorm:"PRIMARY KEY"`
	AdminName string `gorm:"NOT NULL"`
	Function  string `gorm:"DEFAULT '管理日常事务'"`
}

func GetAdministratorById(ADId string) *Administrator {
	var admin Administrator
	result := GlobalConn.Where(&Administrator{AdminId: ADId}).First(&admin)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &admin
}

func GetAllAdministrators() []Administrator {
	var admins = make([]Administrator, 0)
	result := GlobalConn.Model(&Administrator{}).Find(&admins)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return admins
}

func AdministratorsToViewUser(admins []Administrator) []ViewUser {
	viewUsers := make([]ViewUser, 0)
	users := GetAllUser()
	SignInUsers := make(map[string]bool)
	for _, user := range users {
		SignInUsers[user.Id] = true
	}
	for _, admin := range admins {
		viewUsers = append(viewUsers, ViewUser{
			Key:       admin.AdminId,
			UserId:    admin.AdminId,
			UserName:  admin.AdminName,
			Identity:  "administrators",
			HasSignIn: SignInUsers[admin.AdminId],
		})
	}
	return viewUsers
}

func (a *Administrator) AddAdministrator() error {
	var admin Administrator
	result := GlobalConn.Where(&Administrator{AdminId: a.AdminId}).First(&admin)
	if result.RowsAffected != 0 {
		fmt.Println("该管理员已存在")
		return errors.New("管理员已存在")
	}
	GlobalConn.Create(&a)
	return nil
}
