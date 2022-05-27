package main

import (
	"StuManager/Controller"
	"StuManager/Model"
	"fmt"
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

func setup() {
	Model.OpenDatabase()
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

func TestViewAllCourseRouter(t *testing.T) {
	setup()
	router := gin.Default()
	router.Use(TestController.Cors())
	router.POST("/viewAllCourse", TestController.ViewAllCourse)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/viewAllCourse", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	fmt.Println(w.Body.String())
	tearDown()
}
