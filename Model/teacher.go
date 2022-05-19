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
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &t
}

func GetTeacherByName(TName string) *Teacher {
	var t Teacher
	result := GlobalConn.Model(&Teacher{}).Where("teacher_name=?", TName).First(&t)
	if result.Error != nil || result.RowsAffected == 0 {
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
	GlobalConn.Create(&teacher)
	return nil
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
