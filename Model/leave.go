package Model

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"time"
)

type LeaveList struct {
	LeaveID              string `gorm:"PRIMARY KEY"`
	CourseID             string `gorm:"NOT NULL" json:"leaveCourseId"`
	LeaveTime            string `gorm:"NOT NULL" json:"leaveTime"`
	ApplicateTime        string `gorm:"NOT NULL" `
	Reason               string `gorm:"NOT NULL" json:"leaveReason"`
	ApplicantID          string `gorm:"NOT NULL" json:"userId"`
	ReviewerID           string `gorm:"NOT NULL"`
	TeacherOpinion       string `gorm:"DEFAULT:'pending'" json:"teacherIdea"`
	AdministratorOpinion string `gorm:"DEFAULT:'pending'" json:"administratorIdea"`
	Result               string `gorm:"DEFAULT:'pending'"`
}

type ViewLeave struct {
	Key                  string `json:"key"`
	LeaveID              string `json:"leaveId"`
	LeaveCourseId        string `json:"leaveCourseId"`
	LeaveCourseName      string `json:"leaveCourseName"`
	LeaveTime            string `json:"leaveTime"`
	LeaveReason          string `json:"leaveReason"`
	TeacherName          string `json:"teacher"`
	ApplicantId          string `json:"applicantId"`
	ApplicantName        string `json:"applicantName"`
	Address              string `json:"place"`
	College              string `json:"courseCollege"`
	TeacherOpinion       string `gorm:"DEFAULT:'pending'" json:"teacherIdea"`
	AdministratorOpinion string `gorm:"DEFAULT:'pending'" json:"administratorIdea"`
	Result               string `json:"hasPassed"`
}

func (leave *LeaveList) AddLeave() {
	leave.LeaveID = strconv.Itoa(int(crc32.ChecksumIEEE([]byte(leave.CourseID + leave.ApplicantID + time.Now().String()))))
	leave.ApplicateTime = time.Now().Format("2006-01-02 15:04:05")
	leave.TeacherOpinion = "pending"
	leave.AdministratorOpinion = "pending"
	GlobalConn.Create(leave)
}

func (leave *LeaveList) ChangeAdministratorStatus(status string, ReviewerID string) {
	if status == "pass" {
		leave.AdministratorOpinion = "true"
	} else if status == "refuse" {
		leave.AdministratorOpinion = "false"
	} else {
		return
	}
	GlobalConn.Model(&LeaveList{}).Where(&LeaveList{LeaveID: leave.LeaveID}).Update("administrator_opinion", leave.AdministratorOpinion)
	GlobalConn.Model(&LeaveList{}).Where(&LeaveList{LeaveID: leave.LeaveID}).Update("reviewer_id", ReviewerID)
}

func (leave *LeaveList) ChangeTeacherStatus(status string) {
	if status == "pass" {
		leave.TeacherOpinion = "true"
		fmt.Println("审核通过")
		result := GlobalConn.Model(&LeaveList{}).Where(&LeaveList{LeaveID: leave.LeaveID}).Update("teacher_opinion", leave.TeacherOpinion)
		fmt.Println(result.Error)
	} else if status == "refuse" {
		leave.TeacherOpinion = "false"
	} else {
		return
	}

}

func GetLeaveByUserId(UId string) []LeaveList {
	var leaves []LeaveList
	GlobalConn.Where(&LeaveList{ApplicantID: UId}).Find(&leaves)
	return leaves
}

func GetAllLeave() []LeaveList {
	var leaves []LeaveList
	GlobalConn.Find(&leaves)
	return leaves
}

func GetLeaveByLeaveId(LId string) LeaveList {
	var leave LeaveList
	GlobalConn.Where(&LeaveList{LeaveID: LId}).First(&leave)
	return leave
}

func (leave *LeaveList) GetViewLeave() *ViewLeave {
	course := GetCourseById(leave.CourseID)
	college := GetCollegeByNum(course.College)
	timetable := Timetable{}
	timetable.CopyCourse(*course)
	//fmt.Println(timetable.Teacher)
	user := GetUserById(leave.ApplicantID)
	view := ViewLeave{
		Key:                  leave.LeaveID,
		LeaveID:              leave.LeaveID,
		LeaveCourseId:        timetable.CourseId,
		LeaveCourseName:      timetable.CourseName,
		LeaveTime:            leave.LeaveTime,
		LeaveReason:          leave.Reason,
		TeacherName:          timetable.Teacher,
		ApplicantId:          leave.ApplicantID,
		ApplicantName:        user.Name,
		Address:              timetable.Address,
		College:              college.CollegeName,
		TeacherOpinion:       leave.TeacherOpinion,
		AdministratorOpinion: leave.AdministratorOpinion,
		Result:               leave.Result,
	}
	return &view
}
