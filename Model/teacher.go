package Model

import (
	"errors"
	"fmt"
)

type Teacher struct {
	TeacherID   string `gorm:"PRIMARY KEY"`
	TeacherName string `gorm:"NOT NULL"`
	Gender      string `gorm:"NOT NULL"`
	Title       string `gorm:"DEFAULT '无'"`
	College     string `gorm:"NOT NULL"`
}

func NewTeacher() *Teacher {
	return &Teacher{}
}

func GetTeacherById(TId string) *Teacher {
	var t Teacher
	result := GlobalConn.Where(&Teacher{TeacherID: TId}).First(&t)
	if result.Error != nil {
		return nil
	}
	return &t
}

func GetTeacherByName(TName string) *Teacher {
	var t Teacher
	result := GlobalConn.Model(&Teacher{}).Where("teacher_name=?", TName).First(&t)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil
	}
	return &t
}

func (teacher *Teacher) AddTeacher() error {
	var t Teacher
	result := GlobalConn.Where(&Teacher{TeacherID: teacher.TeacherID}).First(&t)
	if result.RowsAffected != 0 {
		fmt.Println("该教师已存在")
		return errors.New("该教师已存在")
	}
	if teacher.Gender != "男" && teacher.Gender != "女" {
		return errors.New("性别不存在")
	} else if !GlobalTitle[teacher.Title] {
		return errors.New("职称不存在")
	} else if GetCollegeByName(teacher.College) == nil {
		return errors.New("学院不存在")
	}
	GlobalConn.Create(&teacher)
	return nil
}

func GetAllTeachers() []Teacher {
	var ters = make([]Teacher, 0)
	result := GlobalConn.Model(&Teacher{}).Find(&ters)
	if result.Error != nil{
		return nil
	}
	return ters
}

func TeachersToViewUser(teachers []Teacher) []ViewUser {
	viewUsers := make([]ViewUser, 0)
	users := GetAllUser()
	SignInUsers := make(map[string]bool)
	for _, user := range users {
		SignInUsers[user.Id] = true
	}
	for _, teacher := range teachers {
		viewUsers = append(viewUsers, ViewUser{
			Key:       teacher.TeacherID,
			UserId:    teacher.TeacherID,
			UserName:  teacher.TeacherName,
			Identity:  "teacher",
			HasSignIn: SignInUsers[teacher.TeacherID],
		})
	}
	return viewUsers
}

var GlobalTeacher = []Teacher{
	{
		TeacherID:   "12",
		TeacherName: "赵宏宇",
		Gender:      "男",
		Title:       "副教授",
		College:     "计算机与人工智能学院",
	},
	{
		TeacherID:   "22",
		TeacherName: "陈剑波",
		Gender:      "男",
		Title:       "讲师",
		College:     "计算机与人工智能学院",
	},
	{
		TeacherID:   "32",
		TeacherName: "何太军",
		Gender:      "男",
		Title:       "副教授",
		College:     "计算机与人工智能学院",
	},
	{
		TeacherID:   "42",
		TeacherName: "黄海于",
		Gender:      "男",
		Title:       "副教授",
		College:     "计算机与人工智能学院",
	},
	{
		TeacherID:   "52",
		TeacherName: "马淑霞",
		Gender:      "女",
		Title:       "副教授",
		College:     "数学学院",
	},
	{
		TeacherID:   "2",
		TeacherName: "崔晓东",
		Gender:      "男",
		Title:       "讲师",
		College:     "马克思主义学院",
	},
}

func TestAddTeacher() {
	GlobalConn.DropTableIfExists(&Teacher{})
	GlobalConn.AutoMigrate(&Teacher{})
	for _, t := range GlobalTeacher {
		t.AddTeacher()
	}
}

var GlobalTitle = map[string]bool{
	"讲师":  true,
	"副教授": true,
	"教授":  true,
}
