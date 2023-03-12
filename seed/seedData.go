package main

import (
	"fmt"
	"schoolapi/initialisers"
	"schoolapi/models"

	"gorm.io/gorm"
)

func init() {
	initialisers.LoadEnvVariables()
	initialisers.ConnectToDB()
}

func main() {
	teacherKen := models.Teacher{Email: "teacherken@gmail.com"}
	teacherJoe := models.Teacher{Email: "teacherjoe@gmail.com"}

	studentsUnderKen := []models.Student{
		{Email: "studentjon@gmail.com", Suspended: false},
		{Email: "studenthon@gmail.com", Suspended: false},
		{Email: "student_only_under_teacher_ken@gmail.com", Suspended: false},
	}

	studentsUnderJoe := []models.Student{
		{Email: "studentmary@gmail.com", Suspended: false},
		{Email: "studentbob@gmail.com", Suspended: false},
		{Email: "studentagnes@gmail.com", Suspended: false},
		{Email: "studentmiche@gmail.com", Suspended: false},
	}

	studentsUnderBoth := []models.Student{
		{Email: "commonstudent1@gmail.com", Suspended: false},
		{Email: "commonstudent2@gmail.com", Suspended: false},
	}

	initialisers.DB.Transaction(func(tx *gorm.DB) error {
		var result *gorm.DB

		// Insert the two teachers first
		result = tx.Create(&teacherKen)
		if !didCreationSucceed(result.Error, teacherKen.Email) {
			return result.Error
		}

		result = tx.Create(&teacherJoe)
		if !didCreationSucceed(result.Error, teacherKen.Email) {
			return result.Error
		}

		for _, student := range studentsUnderKen {
			result = tx.Create(&student)
			if !didCreationSucceed(result.Error, student.Email) {
				return result.Error
			}

			tx.Model(&teacherKen).Association("Students").Append(&student)
			tx.Model(&student).Association("Teachers").Append(&teacherKen)
		}

		for _, student := range studentsUnderJoe {
			result = tx.Create(&student)
			if !didCreationSucceed(result.Error, student.Email) {
				return result.Error
			}

			tx.Model(&teacherJoe).Association("Students").Append(&student)
			tx.Model(&student).Association("Teachers").Append(&teacherJoe)
		}

		for _, student := range studentsUnderBoth {
			result = tx.Create(&student)
			if !didCreationSucceed(result.Error, student.Email) {
				return result.Error
			}

			tx.Model(&teacherKen).Association("Students").Append(&student)
			tx.Model(&student).Association("Teachers").Append(&teacherKen)
			tx.Model(&teacherJoe).Association("Students").Append(&student)
			tx.Model(&student).Association("Teachers").Append(&teacherJoe)
		}

		fmt.Println("Successfully seeded database.")

		// Commit transaction
		return nil
	})
}

func didCreationSucceed(creationResult error, identifier string) bool {
	if creationResult != nil {
		fmt.Println("Error on insertion of '" + identifier + "'. Rolling back transaction")
		return false
	}
	return true
}
