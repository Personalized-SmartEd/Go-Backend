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

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classroom)

	}
}

func GetStudentsByClassroomID() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classroom.Students)

	}
}

func GetTeachersByClassroomID() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, classroom.TeacherID)

	}
}

func GetClassroomsBySchoolCode() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
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

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var classroom utils.InputClassroom
		if err := c.BindJSON(&classroom); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var inputClassroom models.Classroom
		inputClassroom.TeacherID = classroom.TeacherID
		inputClassroom.SchoolCode = classroom.SchoolCode
		inputClassroom.ClassNumber = classroom.ClassNumber
		inputClassroom.ClassCode = classroom.ClassCode

		inputClassroom.ID = primitive.NewObjectID()

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

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")
		studentID := c.Param("student_id")

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updatedClassroom, updateErr := classroomCollection.UpdateOne(ctx, bson.M{"_id": classroomID}, bson.M{"$push": bson.M{"students": studentID}})
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while updating classroom"})
			return
		}

		c.JSON(http.StatusOK, updatedClassroom)

	}
}

func LeaveClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")
		studentID := c.Param("student_id")

		var classroom models.Classroom
		err := classroomCollection.FindOne(ctx, bson.M{"_id": classroomID}).Decode(&classroom)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updatedClassroom, updateErr := classroomCollection.UpdateOne(ctx, bson.M{"_id": classroomID}, bson.M{"$pull": bson.M{"students": studentID}})
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while updating classroom"})
			return
		}

		c.JSON(http.StatusOK, updatedClassroom)

	}
}

func UpdateClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		var classroom models.Classroom
		if err := c.BindJSON(&classroom); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updatedClassroom, err := classroomCollection.UpdateOne(ctx, bson.M{"classroom_id": classroomID}, bson.M{"$set": classroom})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedClassroom)

	}
}

func DeleteClassroom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		classroomID := c.Param("classroom_id")

		_, err := classroomCollection.DeleteOne(ctx, bson.M{"_id": classroomID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "Classroom deleted")

	}
}
