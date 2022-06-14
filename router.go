package main

import (
	"StuManager/Controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	controller := Controller.NewController()
	router.Use(controller.Cors())
	//router.Use(Model.JwtVerfiy)

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
	router.POST("/changePassword", controller.ChangePassword)
	router.POST("/commitLeave", controller.CommitLeave)
	router.POST("/deleteCourse", controller.DeleteCourse)
	router.POST("/deleteSelectedCourse", controller.DeleteSelectedCourse)
	router.POST("/delUser", controller.DeleteUser)
	router.POST("/examLeave", controller.ExamLeave)
	router.POST("/examLeaveByTeacher", controller.ExamLeaveByTeacher)
	router.POST("/getQuestion", controller.GetQuestion)
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
	router.POST("/viewAllStuInACourse", controller.ViewAllStuInACourse)
	router.POST("/viewAllTeachCourse", controller.ViewAllTeachCourse)
	router.POST("/viewAlreadyRegisteredUsers", controller.ViewAlreadyRegisteredUsers)
	router.POST("/viewInitialCourse", controller.ViewInitialCourse)
	router.POST("/viewMyLeave", controller.ViewMyLeave)
	router.POST("/viewSelectedCourse", controller.ViewSelectedCourse)
	router.POST("/viewSelectedCoursesByWeek", controller.ViewSelectedCoursesByWeek)
	router.POST("/viewStuLeave", controller.ViewStuLeave)
	router.POST("/viewStuLeaveByTeacher", controller.ViewStuLeaveByTeacher)
	router.POST("/viewUser", controller.ViewAllUser)
	router.POST("/whetherTeaching", controller.WhetherTeaching)
	router.POST("/autoLogin", controller.AutoLogin)
	router.POST("/exitLogin", controller.ExitLogin)
	return router
}
