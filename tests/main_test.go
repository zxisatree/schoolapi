package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"schoolapi/controllers"
	"schoolapi/initialisers"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	initialisers.LoadEnvVariables()
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestCommonStudents(t *testing.T) {
	initialisers.InitialiseMockDb(t)
	defer initialisers.MockDB.Close()

	initialisers.Mock.ExpectQuery(regexp.QuoteMeta(`as relevant_teachers FROM`)).
		WillReturnRows(sqlmock.NewRows([]string{"email", "relevant_teachers"}).
			AddRow("studentmary@gmail.com", 2).
			AddRow("studentbob@gmail.com", 2).
			AddRow("commonstudent1@gmail.com", 2).
			AddRow("commonstudent2@gmail.com", 2))
	mockResponse := `{"students":["studentmary@gmail.com","studentbob@gmail.com","commonstudent1@gmail.com","commonstudent2@gmail.com"]}`

	r := SetUpRouter()
	r.GET("/api/commonstudents", controllers.GetCommonStudents)
	req, err := http.NewRequest("GET", `http://127.0.0.1:`+os.Getenv("PORT")+`/api/commonstudents?teacher=teacherken%40gmail.com&teacher=teacherjoe%40gmail.com`, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterStudents(t *testing.T) {
	initialisers.InitialiseMockDb(t)
	defer initialisers.MockDB.Close()

	initialisers.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "suspended"}).
			AddRow(1, time.Now(), time.Now(), nil, "studentmary@gmail.com", 0))
	initialisers.Mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	initialisers.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := SetUpRouter()
	r.POST("/api/register", controllers.RegisterStudents)
	req, err := http.NewRequest("POST", `http://127.0.0.1:`+os.Getenv("PORT")+`/api/register`, bytes.NewBuffer([]byte(`{
		"teacher": "teacherken@gmail.com",
		"students":
		[
			"studentmary@gmail.com"
		]
	}`)))
	if err != nil {
		fmt.Println(err)
		return
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestSuspendStudent(t *testing.T) {
	initialisers.InitialiseMockDb(t)
	defer initialisers.MockDB.Close()

	initialisers.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "suspended"}).
			AddRow(3, time.Now(), time.Now(), nil, "studentmary@gmail.com", 1))
	initialisers.Mock.ExpectBegin()
	initialisers.Mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	initialisers.Mock.ExpectCommit()

	r := SetUpRouter()
	r.POST("/api/suspend", controllers.SuspendStudent)
	req, err := http.NewRequest("POST", `http://127.0.0.1:`+os.Getenv("PORT")+`/api/suspend`, bytes.NewBuffer([]byte(`{
		"student": "studentmary@gmail.com"
	}`)))
	if err != nil {
		fmt.Println(err)
		return
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestGetNotifiableStudents(t *testing.T) {
	initialisers.InitialiseMockDb(t)
	defer initialisers.MockDB.Close()

	initialisers.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "email", "suspended"}).
			AddRow(1, time.Now(), time.Now(), nil, "studentbob@gmail.com", 0).
			AddRow(2, time.Now(), time.Now(), nil, "studentmiche@gmail.com", 0))

	mockResponse := `{"recipients":["studentbob@gmail.com","studentmiche@gmail.com"]}`

	r := SetUpRouter()
	r.POST("/api/retrievefornotifications", controllers.GetNotifiableStudents)
	req, err := http.NewRequest("POST", `http://127.0.0.1:`+os.Getenv("PORT")+`/api/retrievefornotifications`, bytes.NewBuffer([]byte(`{
		"teacher": "teacherken@gmail.com",
		"notification": "Hello students! @studentmiche@gmail.com"
	}`)))
	if err != nil {
		fmt.Println(err)
		return
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}
