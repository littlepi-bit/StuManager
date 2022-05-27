package main

import (
	"StuManager/Model"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	Model.OpenDatabase()
	defer Model.CloseDatabase()
	//Model.InitDatabase()
	//Model.TestAddCollege()
	//Model.TestAddTeacher()
	//Model.TestAddCourse()
	//Model.TestAddStudent()
	//Model.InitUserTable()
	//Model.GlobalConn.CreateTable(&Model.LeaveList{})
	//Model.GlobalConn.Create(&Model.Administrator{AdminId: "001", AdminName: "王一博"})
	router := SetUpRouter()
	router.Run(":8000")
}
