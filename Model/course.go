package Model

import (
	"fmt"
	"hash/crc32"
	"log"
	"strconv"
	"sync"
)

type Course struct {
	CourseId    string `gorm:"primaryKey autoIncrement" json:"courseId"`
	CourseName  string `gorm:"NOT NULL" json:"courseName"`
	CourseTime  string `gorm:"NOT NULL" json:"courseTime"`
	Credit      int    `gorm:"NOT NULL" json:"courseValue"`
	Capacity    int    `gorm:"DEFAULT 0" json:"maxCapacity"`
	NumberOfStu int    `gorm:"DEFAULT 0" `
	TeacherID   string `gorm:"DEFAULT '无'" json:"teacher"`
	College     string `gorm:"DEFAULT '无'" json:"courseCollege"`
	Address     string `gorm:"DEFAULT '地点待定'" json:"place"`
	Agreed      string `json:"hasAgreed"`
}

type Timetable struct {
	Key         string `json:"key"`
	CourseId    string `json:"courseId"`
	CourseName  string `json:"courseName"`
	CourseTime  string `json:"courseTime"`
	College     string `json:"courseCollege"`
	Credit      string `json:"courseValue"`
	Capacity    string `json:"capacity"`
	MaxCapacity string `json:"maxCapacity"`
	HasAgreed   string `json:"hasAgreed"`
	Teacher     string `json:"teacher"`
	Address     string `json:"place"`
	HasSelected bool   `json:"hasSelected"`
}

type WeekTable struct {
	Key        string `json:"key"`
	Week       int    `json:"week"`
	Begin      int    `json:"begin"`
	End        int    `json:"end"`
	CourseId   string `json:"courseId"`
	CourseName string `json:"courseName"`
	Teacher    string `json:"teacher"`
	Address    string `json:"place"`
}

var mu sync.Mutex

func GetCourse(num int) (courses []Course) {
	GlobalConn.Limit(num).Find(&courses)
	return courses
}

func GetAllCourse() (courses []Course) {
	GlobalConn.Find(&courses)
	return courses
}

func GetCourseCount() int64 {
	var count int64
	GlobalConn.Model(&Course{}).Count(&count)
	return count
}

func GetCourseById(CId string) *Course {
	var course Course
	GlobalConn.Where(&Course{CourseId: CId}).Find(&course)
	return &course
}

func GetCourseByName(CName string) (courses []Course) {
	GlobalConn.Where("CourseName = ?", CName).Find(&courses)
	return courses
}

func GetAllTimetable(StuId string) []Timetable {
	courses := GetAllCourse()
	var timetable = []Timetable{}
	for _, course := range courses {
		if course.Agreed != "true" {
			continue
		}
		var class = Timetable{}
		class.CopyCourse(course)
		result := GetSelectionsByCourseIdAndStuId(StuId, course.CourseId)
		if result != nil {
			class.HasSelected = true
		} else {
			class.HasSelected = false
		}
		timetable = append(timetable, class)
	}
	return timetable
}

func (t *Timetable) CopyCourse(course Course) {
	var teacher = Teacher{}
	GlobalConn.Where(&Teacher{TeacherID: course.TeacherID}).First(&teacher)
	college := GetCollegeByNum(course.College)
	t.Key = course.CourseId
	t.Teacher = teacher.TeacherName
	t.CourseId = course.CourseId
	t.CourseName = course.CourseName
	t.CourseTime = course.CourseTime
	t.College = college.CollegeName
	t.Address = course.Address
	t.Capacity = strconv.Itoa(course.NumberOfStu) + "/" + strconv.Itoa(course.Capacity)
	t.Credit = strconv.Itoa(course.Credit)
	t.HasAgreed = course.Agreed
}

func (course *Course) CopyTimetable(t Timetable) {
	fmt.Println(t.Teacher)
	var teacher = GetTeacherById(t.Teacher)
	college := GetCollegeByName(t.College)
	course.CourseId = t.CourseId
	course.CourseName = t.CourseName
	course.CourseTime = t.CourseTime
	course.Address = t.Address
	course.Capacity, _ = strconv.Atoi(t.MaxCapacity)
	course.Credit, _ = strconv.Atoi(t.Credit)
	course.College = college.CollegeNumber
	course.Agreed = t.HasAgreed
	course.TeacherID = teacher.TeacherID
}

func (c *Course) GetCId() {
	c.CourseId = c.College + strconv.Itoa(int(crc32.ChecksumIEEE([]byte(c.CourseId))))
}

func (c *Course) AddCourse() {
	if len(c.CourseId) < 3 {
		c.GetCId()
	}
	GlobalConn.Create(c)
}

func (c *Course) DeleteCourse() {
	GlobalConn.Where(&Course{CourseId: c.CourseId}).Delete(&c)
}

func (c *Course) AddStuNum() {
	mu.Lock()
	if c.NumberOfStu >= c.Capacity {
		log.Println("人数超过上限")
		return
	}
	c.NumberOfStu++
	GlobalConn.Model(&Course{}).Where(&Course{CourseId: c.CourseId}).Update("number_of_stu", c.NumberOfStu)
	mu.Unlock()
}

func (c *Course) SubStuNum() {
	mu.Lock()
	if c.NumberOfStu >= c.Capacity {
		log.Fatal("人数超过上限")
		return
	}
	c.NumberOfStu--
	GlobalConn.Model(&Course{}).Where(&Course{CourseId: c.CourseId}).Update("number_of_stu", c.NumberOfStu)
	mu.Unlock()
}

func (c *Course) ChangeCourseAgreed(attitude string) {
	if attitude == "pass" {
		c.Agreed = "true"
	} else if attitude == "reject" {
		c.Agreed = "false"
	} else {
		return
	}
	GlobalConn.Model((&Course{})).Where(&Course{CourseId: c.CourseId}).Update("agreed", c.Agreed)
}

var GlobalCourse []Course = []Course{
	{
		CourseId:   "1",
		CourseName: "离散数学",
		CourseTime: "周四第二讲",
		Credit:     2,
		Capacity:   100,
		TeacherID:  "12",
		College:    "MATH",
		Address:    "x1416",
		Agreed:     "true",
	},
	{
		CourseId:   "2",
		CourseName: "软件设计实现",
		CourseTime: "周五第二讲",
		Credit:     4,
		Capacity:   80,
		TeacherID:  "22",
		College:    "SCAI",
		Address:    "x2416",
		Agreed:     "true",
	},
	{
		CourseId:   "3",
		CourseName: "计算机图形学",
		CourseTime: "周一第四讲",
		Credit:     5,
		Capacity:   94,
		TeacherID:  "32",
		College:    "SCAI",
		Address:    "x4151",
		Agreed:     "true",
	},
	{
		CourseId:    "4",
		CourseName:  "数据结构",
		CourseTime:  "周三第四讲",
		Credit:      4,
		Capacity:    80,
		NumberOfStu: 80,
		TeacherID:   "12",
		College:     "SCAI",
		Address:     "x4154",
		Agreed:      "true",
	},
	{
		CourseId:   "5",
		CourseName: "动漫与游戏",
		CourseTime: "周二第二讲",
		Credit:     3,
		Capacity:   10,
		TeacherID:  "42",
		College:    "SCAI",
		Address:    "x2412",
		Agreed:     "true",
	},
	{
		CourseId:   "6",
		CourseName: "毛泽东思想与中国特色社会主义理论体系概论",
		CourseTime: "周四第二讲",
		Credit:     2,
		Capacity:   110,
		TeacherID:  "2",
		College:    "MARX",
		Address:    "x2514",
		Agreed:     "true",
	},
}

func TestAddCourse() {
	GlobalConn.DropTableIfExists(&Course{})
	GlobalConn.CreateTable(&Course{})
	for _, c := range GlobalCourse {
		c.AddCourse()
	}
}
