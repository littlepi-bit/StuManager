package Model

import (
	"log"
	"sync"

	"github.com/jinzhu/gorm"
)

// 创建全局连接池句柄
var GlobalConn *gorm.DB
var GBmu sync.Mutex

func OpenDatabase() {
	conn, err := gorm.Open("mysql", "root:123456@(120.77.12.35:3306)/stu?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("failed to connect database")
		return
	}
	GlobalConn = conn
}

func CloseDatabase() {
	GlobalConn.Close()
}

func InitDatabase() {
	//创建课程表
	GlobalConn.CreateTable(&User{})
	GlobalConn.CreateTable(&Course{})
	GlobalConn.CreateTable(&Student{})
	GlobalConn.CreateTable(&Teacher{})
	GlobalConn.CreateTable(&Administrator{})
	GlobalConn.CreateTable(&LeaveList{})
	GlobalConn.CreateTable(&Message{})
	GlobalConn.CreateTable(&Selection{})
	TestAddCollege()
	TestAddCourse()
}

func InitUserTable() {
	GlobalConn.Begin()
	if GlobalConn.HasTable(&User{}) {
		GlobalConn.DropTable(&User{})
	}
	GlobalConn.CreateTable(&User{})
	GlobalConn.Commit()
}
