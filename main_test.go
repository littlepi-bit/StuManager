package main

import (
	"StuManager/Controller"
	"StuManager/Model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/*Todo
1.测试路由
2.测试数据库
...
*/

type TestCase struct {
	code        int             //状态码
	param       string          //参数
	method      string          //请求类型
	desc        string          //描述
	handler     gin.HandlerFunc //方法
	showBody    bool            //是否展示返回
	errMsg      string          //错误信息
	url         string          //链接
	contentType string
	exit1       interface{}
	exit2       interface{}
}

func NewBufferStruct(h gin.H) io.Reader {
	jsonByte, _ := json.Marshal(h)
	return bytes.NewReader(jsonByte)
}

func JsontoString(h gin.H) string {
	jsonByte, _ := json.Marshal(h)
	return string(jsonByte)
}

func NewBufferString(body string) io.Reader {
	return bytes.NewBufferString(body)
}

func PerformRequest(method, url, contentType string, body string, handler gin.HandlerFunc) (c *gin.Context, r *http.Request, w *httptest.ResponseRecorder) {
	router := gin.Default()
	router.Use(TestController.Cors())
	router.Handle(method, url, handler)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	r = httptest.NewRequest(method, url, NewBufferString(body))
	c.Request = r
	c.Request.Header.Set("Content-Type", contentType)
	var user = Model.User{
		Id:       "2019110502",
		Password: "12345",
		Name:     "游城十代",
	}
	r.Header.Set("token", Model.GenerateToken(&Model.JWTClaims{
		UserID:   user.Id,
		Username: user.Name,
		Password: user.Password}))
	router.ServeHTTP(w, r)
	return
}

func call(t *testing.T, testcases []TestCase) {
	setup()
	for k, v := range testcases {
		_, _, w := PerformRequest(v.method, v.url, v.contentType, v.param, v.handler)
		fmt.Printf("第%d个测试用例：%s\n", k+1, v.desc)
		if v.showBody {
			fmt.Printf("接口返回: %s\n", w.Body.String())
		}

		assert.Equal(t, v.code, w.Code)
		if v.errMsg != "" {
			assert.Equal(t, v.errMsg, w.Body.String())
		}

	}
	tearDown()
}

func setup() {
	Model.OpenDatabase(false)
	gin.SetMode(gin.TestMode)
	fmt.Println("Before All Tests")
}

func tearDown() {
	Model.CloseDatabase()
	fmt.Println("After All Tests")
}

var TestController = Controller.NewController()

func TestMain(m *testing.M) {
	setup()
	fmt.Println("Test begins...")
	code := m.Run()
	tearDown()
	os.Exit(code)
}

// func TestViewAllCourseRouter(t *testing.T) {
// 	setup()
// 	call(t, GlobalTestCases)
// 	// router := gin.Default()
// 	// router.Use(TestController.Cors())
// 	// router.POST("/viewAllCourse", TestController.ViewAllCourse)
// 	// w := httptest.NewRecorder()
// 	// jsonByte, _ := json.Marshal(gin.H{"userId": "2019110502"})
// 	// req, _ := http.NewRequest("POST", "/viewAllCourse", bytes.NewReader(jsonByte))
// 	// var user = Model.User{
// 	// 	Id:       "2019110502",
// 	// 	Password: "12345",
// 	// 	Name:     "游城十代",
// 	// }
// 	// req.Header.Set("token", Model.GenerateToken(&Model.JWTClaims{
// 	// 	UserID:   user.Id,
// 	// 	Username: user.Name,
// 	// 	Password: user.Password}))
// 	// router.ServeHTTP(w, req)

// 	// assert.Equal(t, 200, w.Code)
// 	// fmt.Println(w.Body.String())
// 	tearDown()
// }

// func TestViewAllNeedTeach(t *testing.T) {
// 	setup()
// 	router := gin.Default()
// 	router.Use(TestController.Cors())
// 	router.POST("/viewAllNeedTeach", TestController.ViewAllNeedTeach)
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/viewAllNeedTeach", NewBufferStruct(gin.H{"userId": "12"}))
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, 200, w.Code)
// 	fmt.Println(w.Body.String())
// 	tearDown()
// }

func TestAllRouter(t *testing.T) {
	call(t, GlobalTestCases)
}

var GlobalTestCases = []TestCase{
	{
		code: 200,
		param: JsontoString(gin.H{
			"userId":   "2019110502",
			"password": "12345",
		}),
		method:   "POST",
		desc:     "测试登录验证(登录成功样例)",
		handler:  TestController.LoginCheck,
		showBody: true,
		errMsg: JsontoString(gin.H{
			"token": Model.GenerateToken(&Model.JWTClaims{
				UserID:   "2019110502",
				Password: "12345",
				Username: "游城十代"}),
			"msg":        "ok",
			"peopleType": "student",
			"userName":   "游城十代",
		}),
		url:         "/loginCheck",
		contentType: "application/json",
	},
	{
		code: 200,
		param: JsontoString(gin.H{
			"userId":   "2019110502",
			"password": "123456",
		}),
		method:   "POST",
		desc:     "测试登录验证(登录失败样例)",
		handler:  TestController.LoginCheck,
		showBody: true,
		errMsg: JsontoString(gin.H{
			"msg":        "fail",
			"peopleType": "",
			"userName":   "",
		}),
		url:         "/loginCheck",
		contentType: "application/json",
	},
	{
		code:        200,
		param:       JsontoString(gin.H{"userId": "2019110502"}),
		method:      "POST",
		desc:        "测试查看所有课程",
		handler:     TestController.ViewAllCourse,
		showBody:    true,
		errMsg:      "",
		url:         "/viewAllCourse",
		contentType: "application/json",
	},
	{
		code:        200,
		param:       JsontoString(gin.H{"userId": "12"}),
		method:      "POST",
		desc:        "测试查看教师所有需要教的课程",
		handler:     TestController.ViewAllNeedTeach,
		showBody:    true,
		errMsg:      "",
		url:         "/viewAllNeedTeach",
		contentType: "application/json",
	},
}
