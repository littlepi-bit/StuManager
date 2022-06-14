package Model

import (
	"hash/crc32"
	"strconv"
	"time"
)

type Selection struct {
	SelectID   string `gorm:"primaryKey"`
	SelectTime string `gorm:"NOT NULL"`
	StuID      string `gorm:"NOT NULL"`
	CourseID   string `gorm:"NOT NULL"`
}

func NewSelection(stu *Student, cor *Course) *Selection {
	return &Selection{
		SelectID:   strconv.Itoa(int(crc32.ChecksumIEEE([]byte(cor.CourseId + stu.StuID + time.Now().String())))),
		StuID:      stu.StuID,
		CourseID:   cor.CourseId,
		SelectTime: GetSystemTime(),
	}
}

//添加选课单
func (s *Selection) AddSelection() {
	GlobalConn.Create(s)
}

//撤销选课单
func (s *Selection) RevokeSelection() {
	var selection Selection
	course := GetCourseById(s.CourseID)
	GlobalConn.Where(&Selection{StuID: s.StuID, CourseID: s.CourseID}).First(&selection)
	course.SubStuNum()
	GlobalConn.Delete(&selection)
}

//查看某一课程的选课情况
func GetSelectionsByCourseId(CId string) (selections []Selection) {
	//selections = make([]Selection, 0)
	GlobalConn.Model(&Selection{}).Where("course_id=?", CId).Find(&selections)
	return selections
}

//查看某一学生的选课情况
func GetSelectionsByStuId(SId string) (selections []Selection) {
	GlobalConn.Model(&Selection{}).Where("stu_id=?", SId).Find(&selections)
	return selections
}

func GetSelectionsByCourseIdAndStuId(StuId, CorId string) *Selection {
	var s Selection
	result := GlobalConn.Where("stu_id = ? AND course_id = ?", StuId, CorId).First(&s)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &s
}
