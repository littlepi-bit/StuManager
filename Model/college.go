package Model

import (
	"hash/crc32"
	"strconv"
)

type College struct {
	CollegeId     string `gorm:"primaryKey"`
	CollegeNumber string `gorm:"NOT NULL"`
	CollegeName   string `gorm:"NOT NULL"`
}

func (college *College) AddCollege() {
	college.CollegeId = strconv.Itoa(int(crc32.ChecksumIEEE([]byte(college.CollegeNumber + college.CollegeName))))
	GlobalConn.Create(&college)
}

func GetCollegeByNum(num string) *College {
	var c College
	result := GlobalConn.Where(&College{CollegeNumber: num}).First(&c)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &c
}

func GetCollegeByName(name string) *College {
	var c College
	result := GlobalConn.Where(&College{CollegeName: name}).First(&c)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	return &c
}

var GlobalCollege = []College{
	{
		CollegeNumber: "SIST",
		CollegeName:   "信息科学与技术学院",
	},
	{
		CollegeNumber: "SCAI",
		CollegeName:   "计算机与人工智能学院",
	},
	{
		CollegeNumber: "CIVE",
		CollegeName:   "土木工程学院",
	},
	{
		CollegeNumber: "SoFL",
		CollegeName:   "外国语学院",
	},
	{
		CollegeNumber: "MARX",
		CollegeName:   "马克思主义学院",
	},
	{
		CollegeNumber: "MATH",
		CollegeName:   "数学学院",
	},
	{
		CollegeNumber: "PHYS",
		CollegeName:   "物理科学与技术学院",
	},
	{
		CollegeNumber: "PHYE",
		CollegeName:   "体育学院",
	},
}

func TestAddCollege() {
	GlobalConn.DropTable(&College{})
	GlobalConn.CreateTable(&College{})
	for _, c := range GlobalCollege {
		c.AddCollege()
	}
}
