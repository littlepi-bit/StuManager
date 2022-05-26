package Model

import (
	"errors"
	"fmt"
)

type User struct {
	Id        string `gorm:"PRIMARY KEY" json:"userId"`
	Name      string `gorm:"NOT NULL UNIQUE" json:"userName"`
	Password  string `gorm:"NOT NULL" json:"password"`
	Identity  string `gorm:"DEFAULT 'tourist'" json:"peopleType"`
	UserEmail string `json:"userEmail"`
	Question1 string `json:"question1"`
	Answer1   string `json:"answer1"`
	Question2 string `json:"question2"`
	Answer2   string `json:"answer2"`
}

type ViewUser struct {
	Key       string `json:"key"`
	UserId    string `json:"userId"`
	UserName  string `json:"userName"`
	HasSignIn bool   `json:"hasSignIn"`
	Identity  string `json:"peopleType"`
}

var GlobalUser *User

func (user *User) AddUser() error {
	result := GlobalConn.Where(&User{Id: user.Id}).Find(&user)
	if result.Error == nil || result.RowsAffected != 0 {
		fmt.Println("用户已存在")
		return errors.New("用户已存在")
	}
	GlobalConn.Create(user)
	return nil
}

func (user *User) DeleteUser() error {
	result := GlobalConn.Where(&User{Id: user.Id}).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		fmt.Println("用户不存在")
		return errors.New("用户不存在")
	}
	GlobalConn.Model(&User{}).Where(&User{Id: user.Id}).Delete(&user)
	return nil
}

func GetUserById(UId string) *User {
	var user User
	result := GlobalConn.Where(&User{Id: UId}).First(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &user
}

func GetAllUser() []User {
	var users []User
	result := GlobalConn.Find(&users)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return users
}

func GetAllViewUser() []ViewUser {
	viewUsers := make([]ViewUser, 0)
	stus := GetAllStudents()
	viewUsers = append(viewUsers, StudentsToViewUser(stus)...)
	ters := GetAllTeachers()
	viewUsers = append(viewUsers, TeachersToViewUser(ters)...)
	admins := GetAllAdministrators()
	viewUsers = append(viewUsers, AdministratorsToViewUser(admins)...)
	return viewUsers
}
