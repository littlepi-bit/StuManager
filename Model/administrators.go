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
