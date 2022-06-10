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
	{
		CollegeNumber: "ELEC",
		CollegeName:   "电气工程学院",
	},
	{
		CollegeNumber: "TRAL",
		CollegeName:   "交通运输与物流学院",
	},
	{
		CollegeNumber: "SoEM",
		CollegeName:   "经济管理学院",
	},
	{
		CollegeNumber: "MASE",
		CollegeName:   "材料科学与工程学院",
	},
	{
		CollegeNumber: "FGEE",
		CollegeName:   "地球科学与环境工程学院",
	},
	{
		CollegeNumber: "SoAD",
		CollegeName:   "建筑学院",
	},
	{
		CollegeNumber: "SHUM",
		CollegeName:   "人文学院",
	},
	{
		CollegeNumber: "SPAL",
		CollegeName:   "公共管理学院",
	},
	{
		CollegeNumber: "DESI",
		CollegeName:   "设计艺术学院",
	},
	{
		CollegeNumber: "MECH",
		CollegeName:   "力学与航空航天学院",
	},
	{
		CollegeNumber: "SLSE",
		CollegeName:   "生命科学与工程学院",
	},
	{
		CollegeNumber: "CPRC",
		CollegeName:   "心理研究与咨询中心",
	},
	{
		CollegeNumber: "XJCO",
		CollegeName:   "利兹学院",
	},
	{
		CollegeNumber: "MYHC",
		CollegeName:   "茅以升学院",
	},
	{
		CollegeNumber: "ISCT",
		CollegeName:   "智慧城市与交通学院",
	},
	{
		CollegeNumber: "PAFD",
		CollegeName:   "武装部、军事教研室",
	},
	{
		CollegeNumber: "LIBR",
		CollegeName:   "图书馆",
	},
	{
		CollegeNumber: "GJHZ",
		CollegeName:   "国际合作与交流处（港澳台事务办公室）",
	},
	{
		CollegeNumber: "ENTC",
		CollegeName:   "工程训练中心",
	},
}

func TestAddCollege() {
	GlobalConn.DropTable(&College{})
	GlobalConn.CreateTable(&College{})
	for _, c := range GlobalCollege {
		c.AddCollege()
	}
}
