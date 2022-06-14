package Model

import (
	"errors"
	"fmt"
)

type Student struct {
	StuID     string `gorm:"primaryKey" json:"userId"`
	StuName   string `gorm:"NOT NULL" json:"userName"`
	StuGender string `gorm:"NOT NULL" json:"userSex"`
	Grade     string `gorm:"NOT NULL"`
	College   string `gorm:"NOT NULL"`
	Major     string `gorm:"NOT NULL"`
	Class     string `gorm:"NOT NULL" json:"userClass"`
}

func NewStudent() *Student {
	return &Student{}
}

func GetStudentByID(SId string) *Student {
	var s Student
	result := GlobalConn.Where(&Student{StuID: SId}).First(&s)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &s
}

func GetAllStudents() []Student {
	var stus = make([]Student, 0)
	result := GlobalConn.Model(&Student{}).Find(&stus)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return stus
}

func StudentsToViewUser(students []Student) []ViewUser {
	viewUsers := make([]ViewUser, 0)
	users := GetAllUser()
	SignInUsers := make(map[string]bool)
	for _, user := range users {
		SignInUsers[user.Id] = true
	}
	for _, student := range students {
		if !SignInUsers[student.StuID] {
			continue
		}
		viewUsers = append(viewUsers, ViewUser{
			Key:       student.StuID,
			UserId:    student.StuID,
			UserName:  student.StuName,
			Identity:  "student",
			HasSignIn: SignInUsers[student.StuID],
		})
	}
	return viewUsers
}

func (student *Student) AddStudent() error {
	var s Student
	result := GlobalConn.Where(&Student{StuID: student.StuID}).First(&s)
	if result.RowsAffected != 0 {
		fmt.Println("该学生已存在")
		return errors.New("学生已存在")
	}
	GlobalConn.Create(&student)
	return nil
}

var GlobalStudent = []Student{
	{
		StuID:     "2019110502",
		StuName:   "游城十代",
		StuGender: "男",
		Grade:     "2019级",
		College:   "计算机与人工智能学院",
		Major:     "计算机科学与技术",
		Class:     "计科1班",
	},
	{
		StuID:     "2019110458",
		StuName:   "老熊",
		StuGender: "男",
		Grade:     "2019级",
		College:   "计算机与人工智能学院",
		Major:     "计算机科学与技术",
		Class:     "计科1班",
	},
	{
		StuID:     "2019110463",
		StuName:   "龙哥",
		StuGender: "男",
		Grade:     "2019级",
		College:   "计算机与人工智能学院",
		Major:     "计算机科学与技术",
		Class:     "计科1班",
	},
}

func TestAddStudent() {
	GlobalConn.DropTableIfExists(&Student{})
	GlobalConn.AutoMigrate(&Student{})
	for _, s := range GlobalStudent {
		s.AddStudent()
	}
}
