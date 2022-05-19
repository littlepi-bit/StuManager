package Model

import (
	"errors"
	"fmt"
)

type Student struct {
	StuID     string `gorm:"primaryKey"`
	StuName   string `gorm:"NOT NULL"`
	StuGender string `gorm:"NOT NULL"`
	Grade     string `gorm:"NOT NULL"`
	College   string `gorm:"NOT NULL"`
	Major     string `gorm:"NOT NULL"`
	Class     string `gorm:"NOT NULL"`
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
