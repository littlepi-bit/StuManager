package Controller

import (
	"StuManager/Model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

//查看所有课程
func (controller *Controller) ViewAllCourse(c *gin.Context) {
	//获取token验证用户
	// MyToken := c.GetHeader("token")
	// if MyToken == "" {
	// 	fmt.Println("不合法访问")
	// 	c.JSON(http.StatusForbidden, gin.H{
	// 		"isExist": false,
	// 	})
	// 	return
	// }
	//user := Model.ParseToken(MyToken)
	var user Model.User
	c.Bind(&user)
	if ok := Model.IsExist(user.Id); !ok {
		fmt.Println("用户不存在")
		c.JSON(http.StatusNotFound, gin.H{
			"isExist": false,
		})
		return
	}
	//获取课程表
	var timetable = Model.GetAllTimetable(user.Id)
	c.JSON(http.StatusOK, timetable)
}

//登录验证
func (controller *Controller) LoginCheck(c *gin.Context) {
	user := &Model.User{}
	c.Bind(&user)
	if user.Id == "" {
		log.Println("用户名不能为空")
		c.JSON(http.StatusOK, gin.H{
			"msg": "fail",
		})
		return
	} else if user.Password == "" {
		log.Println("密码不能为空")
		c.JSON(http.StatusOK, gin.H{
			"msg": "fail",
		})
		return
	}
	result := user.CheckUser()
	fmt.Println(user)
	if !result {
		log.Println("密码或用户名错误！")
		c.JSON(http.StatusOK, gin.H{
			"msg": "fail",
		})
		return
	} else {
		token := Model.GenerateToken(&Model.JWTClaims{
			UserID:   user.Id,
			Username: user.Name,
			Password: user.Password})
		Model.SetHash(
			token,
			Model.JsontoString(gin.H{
				"userId":     user.Id,
				"userName":   user.Name,
				"password":   user.Password,
				"peopleType": user.Identity,
			}),
			time.Minute*5)
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "token", //你的cookie的名字
			Value:    token,   //cookie值
			Path:     "/",
			Domain:   "",
			MaxAge:   604800,
			Secure:   false,
			HttpOnly: false,
			// SameSite: 4, //下面是详细解释
		})
		c.JSON(http.StatusOK, gin.H{
			"token":      token,
			"msg":        "ok",
			"peopleType": user.Identity,
			"userName":   user.Name,
		})
	}
}

//单点登录
func (controller *Controller) AutoLogin(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"flag":       false,
			"peopleType": "",
		})
		return
	}
	if cookie == "" {
		log.Println("重新登录")
		c.JSON(http.StatusOK, gin.H{
			"flag":       false,
			"peopleType": "",
		})
		return
	}
	userString, err := Model.GetHash(cookie)
	if err == Model.RedisErr {
		log.Println("用户未登录")
		c.JSON(http.StatusOK, gin.H{
			"flag":       false,
			"peopleType": "",
		})
		return
	}
	var userVerify Model.User
	json.Unmarshal([]byte(userString), &userVerify)
	fmt.Println(userVerify)
	//user := Model.ParseToken(cookie)
	res := userVerify.CheckUser()
	if !res {
		log.Println("用户不存在")
		c.JSON(http.StatusOK, gin.H{
			"flag":       false,
			"peopleType": "",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"flag":       true,
			"peopleType": userVerify.Identity,
			"userId":     userVerify.Id,
			"password":   userVerify.Password,
			"userName":   userVerify.Name,
			"token":      "",
		})
	}
	// user := userVerify
	// var tmp = Model.GetUserById(user.Id)
	// if tmp == nil {
	// 	log.Println("用户不存在")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"flag":       false,
	// 		"peopleType": "",
	// 	})
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{
	// 	"flag":       true,
	// 	"peopleType": tmp.Identity,
	// 	"userId":     tmp.Id,
	// 	"password":   tmp.Password,
	// 	"userName":   tmp.Name,
	// 	"token":      "",
	// })
}

func (controller *Controller) ExitLogin(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token", //你的cookie的名字
		Value:    "",      //cookie值
		Path:     "/",
		Domain:   "",
		MaxAge:   604800,
		Secure:   false,
		HttpOnly: false,
		// SameSite: 4, //下面是详细解释
	})
	c.JSON(http.StatusOK, gin.H{})
}

//更改密码
func (controller *Controller) ChangePassword(c *gin.Context) {
	var changerUser struct {
		UserId      string `json:"userId"`
		Question1   string `json:"question1"`
		Answer1     string `json:"answer1"`
		Question2   string `json:"question2"`
		Answer2     string `json:"answer2"`
		NewPassword string `json:"newPassword"`
	}
	c.Bind(&changerUser)
	user := Model.GetUserById(changerUser.UserId)
	if user == nil {
		log.Println("用户不存在")
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"msg":    "用户不存在",
		})
		return
	}
	if changerUser.Question1 == user.Question1 && changerUser.Answer1 == user.Answer1 &&
		changerUser.Question2 == user.Question2 && changerUser.Answer2 == user.Answer2 {
		user.ChangePassword(changerUser.NewPassword)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"msg":    "更改密码成功",
		})
	} else if changerUser.Question2 == user.Question1 && changerUser.Answer2 == user.Answer1 &&
		changerUser.Question1 == user.Question2 && changerUser.Answer1 == user.Answer2 {
		user.ChangePassword(changerUser.NewPassword)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"msg":    "更改密码成功",
		})
	} else {
		log.Println("验证失败")
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"msg":    "验证失败",
		})
	}
}

//获取问题
func (controller *Controller) GetQuestion(c *gin.Context) {
	var user Model.User
	c.Bind(&user)
	var tmp = Model.GetUserById(user.Id)
	if tmp == nil {
		log.Println("用户不存在")
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"msg":    "用户不存在！请注册",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"msg":       "用户信息获取成功",
		"question1": tmp.Question1,
		"question2": tmp.Question2,
	})

}

//注册
func (controller *Controller) SignIn(c *gin.Context) {
	var myUser Model.User
	c.Bind(&myUser)
	myUser.Id = strings.TrimSpace(myUser.Id)
	myUser.Password = strings.TrimSpace(myUser.Password)
	fmt.Println(myUser)
	var user Model.User
	result := Model.GlobalConn.Where(&Model.User{Id: myUser.Id}).First(&user)
	if result.Error == nil || result.RowsAffected != 0 {
		fmt.Println("Id已存在！")
		c.JSON(http.StatusOK, gin.H{
			"token":  "",
			"msg":    "用户Id重复",
			"status": "fail",
		})
		return
	}
	result = Model.GlobalConn.Where("Name = ?", myUser.Name).First(&user)
	if result.Error == nil || result.RowsAffected != 0 {
		fmt.Println("用户名已存在！")
		c.JSON(http.StatusOK, gin.H{
			"token":  "",
			"msg":    "用户名重复",
			"status": "fail",
		})
		return
	}
	if myUser.Identity == "student" {
		student := Model.GetStudentByID(myUser.Id)
		if student == nil {
			fmt.Println("该学生不存在！")
			c.JSON(http.StatusOK, gin.H{
				"token":  "",
				"msg":    "该学生不存在",
				"status": "fail",
			})
			return
		}
	} else if myUser.Identity == "teacher" {
		teacher := Model.GetTeacherById(myUser.Id)
		if teacher == nil {
			fmt.Println("该教师不存在！")
			c.JSON(http.StatusOK, gin.H{
				"token":  "",
				"msg":    "该教师不存在",
				"status": "fail",
			})
			return
		}
	} else if myUser.Identity == "administrators" {
		admin := Model.GetAdministratorById(myUser.Id)
		if admin == nil {
			fmt.Println("该教师不存在！")
			c.JSON(http.StatusOK, gin.H{
				"token":  "",
				"msg":    "该教师不存在",
				"status": "fail",
			})
			return
		}
	}
	myUser.AddUser()
	c.JSON(http.StatusOK, gin.H{
		"token": Model.GenerateToken(&Model.JWTClaims{
			UserID:   user.Id,
			Username: user.Name,
			Password: user.Password}),
		"msg":    "注册成功",
		"status": "ok",
	})
	fmt.Println("注册成功")
}

//选课
func (controller *Controller) SelectCourse(c *gin.Context) {
	var tmp struct {
		UserId   string `json:"userId"`
		CourseId string `json:"courseId"`
	}
	c.Bind(&tmp)
	var student Model.Student
	result := Model.GlobalConn.Where(&Model.Student{StuID: tmp.UserId}).First(&student)
	if result.RowsAffected == 0 {
		fmt.Println("学生不存在")
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "学生不存在",
		})
		return
	}
	var course Model.Course
	result = Model.GlobalConn.Where(&Model.Course{CourseId: tmp.CourseId}).First(&course)
	if result.RowsAffected == 0 {
		fmt.Println("课程不存在")
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "课程不存在",
		})
		return
	} else if course.Capacity <= course.NumberOfStu {
		fmt.Println("容量不足")
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"message": "容量不足",
		})
		return
	}
	rows, _ := Model.GlobalConn.Model(&Model.Course{}).Select("courses.course_time").Joins("join selections on selections.course_id=courses.course_id").Rows()
	var Result struct {
		CourseTime string
	}
	for rows.Next() {
		Model.GlobalConn.ScanRows(rows, &Result)
		fmt.Println(Result)
		if Result.CourseTime == course.CourseTime {
			fmt.Println("选课冲突")
			c.JSON(http.StatusOK, gin.H{
				"status":  "fail",
				"message": "选课冲突",
			})
			return
		}
	}
	selection := Model.NewSelection(&student, &course)
	selection.AddSelection()
	course.AddStuNum()
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "选课成功",
	})
}

//查看已选课程
func (controller *Controller) ViewSelectedCourse(c *gin.Context) {
	var tmp struct {
		UserId string `json:"userId"`
	}
	c.Bind(&tmp)
	var student Model.Student
	result := Model.GlobalConn.Where(&Model.Student{StuID: tmp.UserId}).First(&student)
	if result.RowsAffected == 0 {
		fmt.Println("学生不存在")
		c.JSON(http.StatusForbidden, gin.H{
			"status": "fail",
			"msg":    "学生不存在",
		})
		return
	}
	var selections []Model.Selection
	Model.GlobalConn.Where(&Model.Selection{StuID: student.StuID}).Find(&selections)
	var timetable []Model.Timetable
	if len(selections) == 0 {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}
	for _, selection := range selections {
		course := Model.GetCourseById(selection.CourseID)
		class := Model.Timetable{}
		class.CopyCourse(*course)
		timetable = append(timetable, class)
	}
	c.JSON(http.StatusOK, timetable)
}

//删除已选课程
func (controller *Controller) DeleteSelectedCourse(c *gin.Context) {
	var tmp struct {
		UserId   string `json:"userId"`
		CourseId string `json:"courseId"`
	}
	c.Bind(&tmp)
	fmt.Println(tmp)
	var student Model.Student
	result := Model.GlobalConn.Where(&Model.Student{StuID: tmp.UserId}).First(&student)
	if result.RowsAffected == 0 {
		fmt.Println("学生不存在")
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "学生不存在",
		})
		return
	}
	var course Model.Course
	result = Model.GlobalConn.Where(&Model.Course{CourseId: tmp.CourseId}).First(&course)
	if result.RowsAffected == 0 {
		fmt.Println("课程不存在")
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "课程不存在",
		})
		return
	} else if course.Capacity <= 0 {
		fmt.Println("容量不足")
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "容量不足",
		})
		return
	}
	Model.GlobalConn.Where(&Model.Selection{StuID: student.StuID, CourseID: course.CourseId}).Delete(&Model.Selection{})
	course.SubStuNum()
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "删除成功",
	})
}

//Cors跨域中间件
func (controller *Controller) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		//c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		//c.Header("Access-Control-Allow-Origin", "http://120.77.12.35:3000")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Token,Id")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func (controller *Controller) SelectAllTeach(c *gin.Context) {

}

//提交请假申请
func (controller *Controller) CommitLeave(c *gin.Context) {
	var leave Model.LeaveList
	c.Bind(&leave)
	if leave.Reason == "" {
		fmt.Println("请假理由不能为空")
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "ok",
			"message": "请假理由不能为空",
		})
		return
	}
	leave.AddLeave()
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "请假成功",
	})
}

//查看请假申请
func (controller *Controller) ViewMyLeave(c *gin.Context) {
	var user Model.User
	c.Bind(&user)
	leaves := Model.GetLeaveByUserId(user.Id)
	var viewLeave = []Model.ViewLeave{}
	for _, leave := range leaves {
		if leave.TeacherOpinion == "pass" {
			leave.TeacherOpinion = "true"
		} else if leave.TeacherOpinion == "refuse" {
			leave.TeacherOpinion = "false"
		}
		view := leave.GetViewLeave()
		viewLeave = append(viewLeave, *view)
	}
	c.JSON(http.StatusOK, viewLeave)
}

//发送邮件
func (controller *Controller) SendMessages(c *gin.Context) {
	var message Model.Message
	c.Bind(&message)
	message.NotifiedID = message.NotifiedID[strings.Index(message.NotifiedID, "(")+1 : strings.Index(message.NotifiedID, ")")]
	u1, u2 := Model.GetUserById(message.NotifierID), Model.GetUserById(message.NotifiedID)
	if u1 == nil {
		fmt.Println("未知用户！")
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"message": "未知用户!",
		})
		return
	}
	if u2 == nil {
		fmt.Println("找不到该用户！")
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"message": "找不到该用户!",
		})
		return
	}
	tmp := Model.NewMessage(u1.Id, u2.Id)
	message.MegID = tmp.MegID
	message.SendTime = time.Now().Format("2006-01-02 15:04:05")
	message.AddMessage()
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "发送成功!",
	})
}

//查看发件箱
func (controller *Controller) ViewAllSended(c *gin.Context) {
	var user Model.User
	c.Bind(&user)
	User := Model.GetUserById(user.Id)
	if User == nil {
		fmt.Println("未知用户！")
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}
	messages := Model.GetSendMessage(user)
	var Sendeds = []Model.SendedMsg{}
	for _, message := range messages {
		sended := message.CopySended()
		Sendeds = append(Sendeds, sended)
	}
	c.JSON(http.StatusOK, Sendeds)
}

//查看收件箱
func (controller *Controller) ViewAllReceived(c *gin.Context) {
	var user Model.User
	c.Bind(&user)
	User := Model.GetUserById(user.Id)
	if User == nil {
		fmt.Println("未知用户！")
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}
	messages := Model.GetRecvMessage(user)
	var Recvs = []Model.RecvMsg{}
	for _, message := range messages {
		recv := message.CopyRecv()
		Recvs = append(Recvs, recv)
	}
	c.JSON(http.StatusOK, Recvs)
}

//阅读邮件
func (controller *Controller) ReadMessage(c *gin.Context) {
	var tmp struct {
		UserId    string `json:"userId"`
		MessageId string `json:"messageID"`
	}
	c.Bind(&tmp)
	user := Model.GetUserById(tmp.UserId)
	if user == nil {
		fmt.Println("用户不存在")
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"message": "用户不存在",
		})
	}
	message := Model.GetMessageById(tmp.MessageId)
	if message == nil {
		fmt.Println("消息不存在")
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"message": "消息不存在",
		})
	}
	message.ReadMessage()
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "记录成功",
	})
}

//查看所有用户
func (controller *Controller) ViewAllUser(c *gin.Context) {
	var tmp struct {
		UserId string `json:"userID"`
	}

	var res = []Model.ViewUser{}
	c.Bind(&tmp)

	result := Model.GetUserById(tmp.UserId)
	if result == nil || result.Identity != "administrators" {
		fmt.Println("管理员不存在")
		c.JSON(http.StatusOK, res)
	}
	// users := Model.GetAllUser()
	// SignInUsers := make(map[string]bool)
	// for _, user := range users {
	// 	res = append(res, Model.ViewUser{
	// 		Key:       user.Id,
	// 		UserId:    user.Id,
	// 		HasSignIn: true,
	// 		Identity:  user.Identity,
	// 	})
	// 	SignInUsers[user.Id] = true
	// }
	//stus := Model.GetAllStudents()
	// for _, stu := range stus {
	// 	if SignInUsers[stu.StuID] {
	// 		continue
	// 	}
	// 	res = append(res, Model.ViewUser{
	// 		Key:       stu.StuID,
	// 		UserId:    stu.StuID,
	// 		HasSignIn: false,
	// 		Identity:  "student",
	// 	})
	// }
	// ters := Model.GetAllTeachers()
	// for _, ter := range ters {
	// 	if SignInUsers[ter.TeacherID] {
	// 		continue
	// 	}
	// 	res = append(res, Model.ViewUser{
	// 		Key:       ter.TeacherID,
	// 		UserId:    ter.TeacherID,
	// 		HasSignIn: false,
	// 		Identity:  "teacher",
	// 	})
	// }
	// admins := Model.GetAllAdministrators()
	// for _, admin := range admins {
	// 	if SignInUsers[admin.AdminId] {
	// 		continue
	// 	}
	// 	res = append(res, Model.ViewUser{
	// 		Key:       admin.AdminId,
	// 		UserId:    admin.AdminId,
	// 		HasSignIn: false,
	// 		Identity:  "administrators",
	// 	})
	// }
	res = Model.GetAllViewUser()
	c.JSON(http.StatusOK, res)
}

//查看学生请假单
func (controller *Controller) ViewStuLeave(c *gin.Context) {
	var user = Model.User{}
	c.Bind(&user)
	user = *Model.GetUserById(user.Id)
	leaves := Model.GetAllLeave()
	viewLeaves := []Model.ViewLeave{}
	if user.Identity == "teacher" {
		t := Model.GetTeacherById(user.Id)
		for _, leave := range leaves {
			viewLeave := *leave.GetViewLeave()
			if viewLeave.TeacherName != t.TeacherName {
				continue
			}
			fmt.Println(viewLeave.TeacherName)
			viewLeaves = append(viewLeaves, viewLeave)
		}
		c.JSON(http.StatusOK, viewLeaves)
		return
	}
	for _, leave := range leaves {
		fmt.Println(leave)
		viewLeaves = append(viewLeaves, *leave.GetViewLeave())
	}
	c.JSON(http.StatusOK, viewLeaves)
}

//添加课程
func (controller *Controller) AddCourse(c *gin.Context) {
	var course Model.Course
	//var t *Model.Timetable
	c.Bind(&course)
	TId := strings.Split(course.TeacherID, "(")
	course.TeacherID = TId[1][:len(TId[1])-1]
	course.CourseId = strconv.Itoa(int(Model.GetCourseCount() + 1))
	r := []rune(course.Address)
	fmt.Println(r[1])
	if r[0] == rune('犀') {
		course.Address = "X" + string(r[2:])
	} else if r[0] == rune('九') {
		course.Address = "J" + string(r[2:])
	} else {
		course.Address = "E" + string(r[2:])
	}
	//course.CopyTimetable(*t)
	college := Model.GetCollegeByName(course.College)
	teacher := Model.GetTeacherById(course.TeacherID)
	var tmp Model.Course
	result := Model.GlobalConn.Where(&Model.Course{CourseId: course.CourseId}).First(&tmp)
	if result.Error == nil || result.RowsAffected != 0 {
		fmt.Println("课程已存在")
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"msg":    "课程已存在",
		})
		return
	}
	course.College = college.CollegeNumber
	course.Agreed = "pending"
	course.TeacherID = teacher.TeacherID
	course.NumberOfStu = 0
	course.GetCId()
	fmt.Println(course)
	course.AddCourse()
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "添加成功",
	})
}

//删除课程
func (controller *Controller) DeleteCourse(c *gin.Context) {
	var tmp struct {
		UserId   string `json:"userId"`
		CourseId string `json:"courseId"`
	}
	c.Bind(&tmp)
	result := Model.GetUserById(tmp.UserId)
	if result == nil || result.Identity != "administrators" {
		fmt.Println("管理员不存在")
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"msg":    "管理员不存在",
		})
		return
	}
	course := Model.GetCourseById(tmp.CourseId)
	course.DeleteCourse()
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "删除成功",
	})
}

//查看已添加的课程
func (controller *Controller) ViewInitialCourse(c *gin.Context) {
	var tmp = []Model.Timetable{}
	courses := Model.GetAllCourse()
	for _, course := range courses {
		var t Model.Timetable
		fmt.Println(course.Agreed)
		t.CopyCourse(course)
		t.MaxCapacity = strconv.Itoa(course.Capacity)
		tmp = append(tmp, t)
	}
	c.JSON(http.StatusOK, tmp)
}

//查看某一门课程中所有的学生
func (controller *Controller) ViewAllStuInACourse(c *gin.Context) {
	var course Model.Course
	c.Bind(&course)
	course = *Model.GetCourseById(course.CourseId)
	if course.Agreed == "false" {
		fmt.Println("教师还未同意")
		c.JSON(http.StatusOK, []Model.Student{})
		return
	}
	type ViewStudent struct {
		Key string `json:"key"`
		Model.Student
	}
	students := make([]ViewStudent, 0)
	selections := Model.GetSelectionsByCourseId(course.CourseId)
	for _, selection := range selections {
		student := Model.GetStudentByID(selection.StuID)
		students = append(students, ViewStudent{
			Key:     student.StuID,
			Student: *student,
		})
	}
	c.JSON(http.StatusOK, students)
}

//审核请假单
func (controller *Controller) ExamLeave(c *gin.Context) {
	var tmp struct {
		LeaveId string `json:"leaveId"`
		UserId  string `json:"userId"`
		Result  string `json:"out"`
	}
	c.Bind(&tmp)
	leave := Model.GetLeaveByLeaveId(tmp.LeaveId)
	leave.ChangeAdministratorStatus(tmp.Result, tmp.UserId)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "审核成功",
	})
}

//添加用户
func (controller *Controller) AddUser(c *gin.Context) {
	var user struct {
		UserId   string `json:"userId"`
		UserName string `json:"userName"`
		Identity string `json:"peopleType"`
		Gender   string `json:"userSex"`
		Grade    string `json:"userGrade"`
		College  string `json:"userCollege"`
		Major    string `json:"userMajor"`
		Class    string `json:"userClass"`
		Title    string `json:"userTitle"`
		Type     string `json:"userType"`
	}
	c.Bind(&user)
	var err error
	if user.Identity == "student" {
		var student = Model.Student{
			StuID:     user.UserId,
			StuName:   user.UserName,
			StuGender: user.Gender,
			Grade:     user.Gender,
			College:   user.College,
			Major:     user.Major,
			Class:     user.Class,
		}
		err = student.AddStudent()
	} else if user.Identity == "teacher" {
		var teacher = Model.Teacher{
			TeacherID:   user.UserId,
			TeacherName: user.UserName,
			Gender:      user.Gender,
			Title:       user.Title,
			College:     user.College,
		}
		err = teacher.AddTeacher()
	} else if user.Identity == "administrators" {
		var a = Model.Administrator{
			AdminId:   user.UserId,
			AdminName: user.UserName,
			Function:  user.Type,
		}
		err = a.AddAdministrator()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"msg":    "用户类别不存在",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "添加成功",
	})
}

//删除用户
func (controller *Controller) DeleteUser(c *gin.Context) {
	var user Model.User
	var tmp struct {
		UserId    string `json:"userId"`
		DelUserId string `json:"delUserId"`
	}
	c.Bind(&tmp)
	admin := Model.GetUserById(tmp.UserId)
	if admin.Identity != "administrators" {
		fmt.Println("权限不足")
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"msg":    "没有管理员权限",
		})
		return
	}
	user = *Model.GetUserById(tmp.DelUserId)
	if err := user.DeleteUser(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"msg":    err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "删除成功",
	})
}

//获取所有教师
func (controller *Controller) GetTeachers(c *gin.Context) {
	var user Model.User
	c.Bind(&user)
	admin := Model.GetAdministratorById(user.Id)
	if admin == nil {
		fmt.Println("管理员不存在")
		c.JSON(http.StatusOK, Model.User{})
		return
	}
	teachers := Model.GetAllTeachers()
	viewUsers := Model.TeachersToViewUser(teachers)
	c.JSON(http.StatusOK, viewUsers)
}

//获取所有注册用户
func (controller *Controller) ViewAlreadyRegisteredUsers(c *gin.Context) {
	users := Model.GetAllUser()
	c.JSON(http.StatusOK, users)
}

//教师

//查看教师所有课程
func (controller *Controller) ViewAllTeachCourse(c *gin.Context) {
	var tmp = []Model.Timetable{}
	var user Model.User
	c.Bind(&user)
	courses := Model.GetAllCourse()
	for _, course := range courses {
		if course.TeacherID != user.Id || course.Agreed != "pending" {
			continue
		}
		var t Model.Timetable
		fmt.Println(course.Agreed)
		t.CopyCourse(course)
		t.MaxCapacity = strconv.Itoa(course.Capacity)
		tmp = append(tmp, t)
	}
	c.JSON(http.StatusOK, tmp)
}

//查看教师需要上的课程
func (controller *Controller) ViewAllNeedTeach(c *gin.Context) {
	var tmp = []Model.Timetable{}
	var user Model.User
	c.Bind(&user)
	courses := Model.GetAllCourse()
	for _, course := range courses {
		if course.TeacherID != user.Id || course.Agreed != "true" {
			continue
		}
		var t Model.Timetable
		fmt.Println(course.Agreed)
		t.CopyCourse(course)
		t.MaxCapacity = strconv.Itoa(course.Capacity)
		tmp = append(tmp, t)
	}
	c.JSON(http.StatusOK, tmp)
}

//查看请假单
func (controller *Controller) ViewStuLeaveByTeacher(c *gin.Context) {
	leaves := Model.GetAllLeave()
	viewLeaves := []Model.ViewLeave{}
	var user Model.User
	c.Bind(&user)
	t := Model.GetTeacherById(user.Id)
	for _, leave := range leaves {
		viewLeave := *leave.GetViewLeave()
		if viewLeave.TeacherName != t.TeacherName {
			continue
		}
		fmt.Println(viewLeave.TeacherName)
		viewLeaves = append(viewLeaves, viewLeave)
	}
	c.JSON(http.StatusOK, viewLeaves)
}

//审核请假单
func (controller *Controller) ExamLeaveByTeacher(c *gin.Context) {
	var tmp struct {
		LeaveId string `json:"leaveId"`
		UserId  string `json:"userId"`
		Result  string `json:"out"`
	}
	c.Bind(&tmp)
	leave := Model.GetLeaveByLeaveId(tmp.LeaveId)
	leave.ChangeTeacherStatus(tmp.Result)
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "审核成功",
	})
}

func (controller *Controller) WhetherTeaching(c *gin.Context) {
	var tmp struct {
		UserId   string `json:"userId"`
		CourseId string `json:"courseID"`
		Attitude string `json:"attitude"`
	}
	c.Bind(&tmp)
	if t := Model.GetTeacherById(tmp.UserId); t == nil {
		c.JSON(http.StatusOK, gin.H{
			"statu": "fail",
			"msg":   "老师不存在",
		})
		return
	}
	course := Model.GetCourseById(tmp.CourseId)
	if course == nil {
		c.JSON(http.StatusOK, gin.H{
			"statu": "fail",
			"msg":   "课程不存在",
		})
		return
	}
	var msg string
	if tmp.Attitude == "pass" {
		msg = "成功选择课程"
	} else if tmp.Attitude == "reject" {
		msg = "成功拒绝课程"
	}
	course.ChangeCourseAgreed(tmp.Attitude)
	c.JSON(http.StatusOK, gin.H{
		"statu": "ok",
		"msg":   msg,
	})
}
