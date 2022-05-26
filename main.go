package main

import (
	"StuManager/Controller"
	"StuManager/Model"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
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
	controller := Controller.NewController()
	router.Use(controller.Cors())
	router.Use(Model.JwtVerfiy)

	router.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, []struct{}{})
	})
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "http://localhost:3000/")
	})
	router.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "http://localhost:3000/")
	})
	router.POST("/addCourse", controller.AddCourse)
	router.POST("/addUser", controller.AddUser)
	router.POST("/changePassword") //todo
	router.POST("/commitLeave", controller.CommitLeave)
	router.POST("/deleteCourse", controller.DeleteCourse)
	router.POST("/deleteSelectedCourse", controller.DeleteSelectedCourse)
	router.POST("/delUser", controller.DeleteUser)
	router.POST("/examLeave", controller.ExamLeave)
	router.POST("/examLeaveByTeacher", controller.ExamLeaveByTeacher)
	router.POST("/getQuestion") //todo
	router.POST("/getTeachers", controller.GetTeachers)
	router.POST("/loginCheck", controller.LoginCheck)
	router.POST("/readMessage", controller.ReadMessage)
	router.POST("/selectCourse", controller.SelectCourse)
	router.POST("/sendMessage", controller.SendMessages)
	router.POST("/signIn", controller.SignIn)
	router.POST("/viewAllCourse", controller.ViewAllCourse)
	router.POST("/viewAllNeedTeach", controller.ViewAllNeedTeach)
	router.POST("/viewAllReceived", controller.ViewAllReceived)
	router.POST("/viewAllSended", controller.ViewAllSended)
	router.POST("/viewAllStuInACourse", controller.ViewAllStuInACourse) //todo
	router.POST("/viewAllTeachCourse", controller.ViewAllTeachCourse)
	router.POST("/viewInitialCourse", controller.ViewInitialCourse)
	router.POST("/viewMyLeave", controller.ViewMyLeave)
	router.POST("/viewSelectedCourse", controller.ViewSelectedCourse)
	router.POST("/viewStuLeave", controller.ViewStuLeave)
	router.POST("/viewStuLeaveByTeacher", controller.ViewStuLeaveByTeacher)
	router.POST("/viewUser", controller.ViewAllUser)
	router.POST("/whetherTeaching", controller.WhetherTeaching)

	router.Run(":8000")
}
