package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var classroomCollection *mongo.Collection = config.OpenCollection(config.Client, "classroom")

func GetClassroomByTeacherID() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		teacherID := c.Param("teacher_id")

		classrooms, err := classroomCollection.Find(ctx, bson.M{"teacher_id": teacherID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classrooms)

	}
}

func GetClassroomByID() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"classroom_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classroom)

	}
}

func GetStudentsByClassroomID() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"classroom_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classroom.Students)

	}
}

func GetTeachersByClassroomID() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"classroom_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classroom.TeacherID)

	}
}

func GetClassroomsBySchoolCode() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		schoolCode := c.Param("school_code")

		classrooms, err := classroomCollection.Find(ctx, bson.M{"school_code": schoolCode})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classrooms)

	}
}

func GetClassroomsBySchoolCodeTeacherID() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		schoolCode := c.Param("school_code")
		teacherID := c.Param("teacher_id")

		classrooms, err := classroomCollection.Find(ctx, bson.M{"school_code": schoolCode, "teacher_id": teacherID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classrooms)

	}
}

func CreateClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		teacherID := c.Param("teacher_id")

		var classroom utils.InputClassroom
		if err := c.BindJSON(&classroom); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var inputClassroom models.Classroom
		inputClassroom.TeacherID = teacherID
		inputClassroom.SchoolCode = classroom.SchoolCode
		inputClassroom.ClassNumber = classroom.ClassNumber
		inputClassroom.ClassCode = classroom.ClassCode

		inputClassroom.ID = primitive.NewObjectID()
		inputClassroom.ClassroomID = inputClassroom.ID.Hex()

		result, err := classroomCollection.InsertOne(ctx, inputClassroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

func JoinClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		studentIDIfc, exists := c.Get("student_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Student ID not found in context"})
			return
		}
		studentID, ok := studentIDIfc.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid student ID format"})
			return
		}

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"classroom_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		flag := true
		for _, id := range classroom.Students {
			if id == studentID {
				flag = false
			}
		}

		if flag {
			classroom.Students = append(classroom.Students, studentID)
		}

		result, updateErr := classroomCollection.UpdateOne(ctx, bson.M{"classroom_id": classroomID}, bson.M{"$set": classroom})

		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while updating classroom"})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

func LeaveClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		studentIDIfc, exists := c.Get("student_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "student ID not found in context"})
			return
		}
		studentID, ok := studentIDIfc.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid student ID format"})
			return
		}

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"classroom_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for i, id := range classroom.Students {
			if id == studentID {
				classroom.Students = append(classroom.Students[:i], classroom.Students[i+1:]...)
				break
			}
		}

		result, updateErr := classroomCollection.UpdateOne(ctx, bson.M{"classroom_id": classroomID}, bson.M{"$set": classroom})

		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while updating classroom"})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

func DeleteClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		_, err := classroomCollection.DeleteOne(ctx, bson.M{"classroom_id": classroomID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "Classroom deleted")

	}
}
