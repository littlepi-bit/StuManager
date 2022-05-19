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
	router.POST("/viewAllCourse", controller.ViewAllCourse)
	router.POST("/loginCheck", controller.LoginCheck)
	router.POST("/signIn", controller.SignIn)
	router.POST("/selectCourse", controller.SelectCourse)
	router.POST("/viewSelectedCourse", controller.ViewSelectedCourse)
	router.POST("/deleteSelectedCourse", controller.DeleteSelectedCourse)
	router.POST("/commitLeave", controller.CommitLeave)
	router.POST("/viewMyLeave", controller.ViewMyLeave)
	router.POST("/sendMessage", controller.SendMessages)
	router.POST("/viewAllSended", controller.ViewAllSended)
	router.POST("/viewAllReceived", controller.ViewAllReceived)
	router.POST("/readMessage", controller.ReadMessage)
	router.POST("/viewUser", controller.ViewAllUser)
	router.POST("/viewInitialCourse", controller.ViewInitialCourse)
	router.POST("/addCourse", controller.AddCourse)
	router.POST("/viewStuLeave", controller.ViewStuLeave)
	router.POST("/examLeave", controller.ExamLeave)
	router.POST("/viewAllTeachCourse", controller.ViewAllTeachCourse)
	router.POST("/viewAllNeedTeach", controller.ViewAllNeedTeach)
	router.POST("/viewStuLeaveByTeacher", controller.ViewStuLeaveByTeacher)
	router.POST("/examLeaveByTeacher", controller.ExamLeaveByTeacher)
	router.POST("/whetherTeaching", controller.WhetherTeaching)
	router.POST("/addUser", controller.AddUser)
	router.POST("delUser", controller.DeleteUser)
	router.Run(":8000")
}
