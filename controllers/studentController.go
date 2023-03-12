package controllers

import (
	"regexp"
	"schoolapi/initialisers"
	"schoolapi/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var mentionedStudentEmailRegex *regexp.Regexp

func init() {
	mentionedStudentEmailRegex = regexp.MustCompile(`@(\S+)`)
}

func GetCommonStudents(c *gin.Context) {
	teachers := c.QueryArray("teacher")

	if len(teachers) == 0 {
		c.JSON(400, gin.H{
			"message": "No teachers specified",
		})
		return
	}

	var students []models.Student
	query := initialisers.DB.Model(&models.Student{}).
		Joins("JOIN teachers_students ON teachers_students.student_id = students.id").
		Joins("JOIN teachers ON teachers_students.teacher_id = teachers.id").
		Where("teachers.email = ?", teachers[0])
	for i, teacher := range teachers {
		if i == 0 {
			continue
		}
		query = query.Or("teachers.email = ?", teacher)
	}

	result := query.Select("students.email, count(teachers.id) as relevant_teachers").
		Group("students.email").
		Having("relevant_teachers = ?", len(teachers)).
		Find(&students)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Invalid teacher email(s) provided.",
		})
		return
	}

	var studentEmails []string
	for _, student := range students {
		studentEmails = append(studentEmails, student.Email)
	}

	c.JSON(200, gin.H{
		"students": studentEmails,
	})
}

func RegisterStudents(c *gin.Context) {
	var body struct {
		Teacher  string   `json:"teacher"`
		Students []string `json:"students"`
	}
	c.Bind(&body)
	var result *gorm.DB

	var teacher models.Teacher
	result = initialisers.DB.Where("email = ?", body.Teacher).First(&teacher)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Invalid teacher email provided.",
		})
		return
	}

	// Add references for both student and teacher to each other
	for _, studentEmail := range body.Students {
		var student models.Student
		initialisers.DB.Where("email = ?", studentEmail).First(&student)
		result := initialisers.DB.Model(&student).Association("Teachers").Append(&teacher)
		if result != nil {
			c.JSON(400, gin.H{
				"message": "Invalid student email provided.",
			})
			return
		}

		result = initialisers.DB.Model(&teacher).Association("Students").Append(&student)
		if result != nil {
			c.JSON(400, gin.H{
				"message": "Invalid student email provided.",
			})
			return
		}

	}

	c.Status(204)
}

func SuspendStudent(c *gin.Context) {
	var body struct {
		Student string `json:"student"`
	}
	c.Bind(&body)
	var result *gorm.DB

	var student models.Student
	result = initialisers.DB.Where("email = ?", body.Student).First(&student)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Invalid student email provided.",
		})
		return
	}

	result = initialisers.DB.Model(&student).Updates(models.Student{
		Suspended: true,
	})

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error, please try again in a few seconds.",
		})
		return
	}

	c.Status(204)
}

func GetNotifiableStudents(c *gin.Context) {
	var body struct {
		Teacher      string `json:"teacher"`
		Notification string `json:"notification"`
	}
	c.Bind(&body)

	// Check if teacher exists
	var result *gorm.DB
	var teacher models.Teacher
	result = initialisers.DB.Where("email = ?", body.Teacher).First(&teacher)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Invalid teacher email provided.",
		})
		return
	}

	mentionedStudentMatches := mentionedStudentEmailRegex.FindAllStringSubmatch(body.Notification, -1)
	var mentionedStudentEmails []string
	for _, match := range mentionedStudentMatches {
		mentionedStudentEmails = append(mentionedStudentEmails, match[1])
	}

	var students []models.Student
	query := initialisers.DB.Model(&models.Student{}).
		Joins("JOIN teachers_students ON teachers_students.student_id = students.id").
		Joins("JOIN teachers ON teachers_students.teacher_id = teachers.id").
		Where("students.suspended = 0").
		Where("teachers.email = ?", body.Teacher)
	for i, mentionedStudentEmail := range mentionedStudentEmails {
		if i == 0 {
			continue
		}
		query = query.Or("students.email = ?", mentionedStudentEmail)
	}
	result = query.Find(&students)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Invalid teacher email and/or notification.",
		})
		return
	}

	var studentEmails []string
	for _, student := range students {
		studentEmails = append(studentEmails, student.Email)
	}

	c.JSON(200, gin.H{
		"recipients": studentEmails,
	})
}
