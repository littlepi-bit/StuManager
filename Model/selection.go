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
		SelectTime: time.Now().Format("2006-01-02 15:04:05"),
	}
}

//添加选课单
func (s *Selection) AddSelection() {
	GlobalConn.Create(s)
}

//撤销选课单
func (s *Selection) RevokeSelection() {
	var selection Selection
	GlobalConn.Where(&Selection{StuID: s.StuID, CourseID: s.CourseID}).First(&selection)
	GlobalConn.Delete(&selection)
}
