package controllers

import (
	"backend/internal/config"
	"backend/internal/helper"
	"backend/internal/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var teacherCollection *mongo.Collection = config.OpenCollection(config.Client, "teacher")

func GetTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		teacher_id := c.Param("teacher_id")

		var teacher models.Teacher

		err := teacherCollection.FindOne(ctx, bson.M{"teacher_id": teacher_id}).Decode(&teacher)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing teacher items"})
			return
		}
		c.JSON(http.StatusOK, teacher)
	}
}

func SignUpTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var teacher models.Teacher

		if err := c.BindJSON(&teacher); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(teacher)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := teacherCollection.CountDocuments(ctx, bson.M{"email": teacher.Email})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := HashPassword(teacher.Password)
		teacher.Password = password

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exsits"})
			return
		}

		teacher.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		teacher.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		teacher.ID = primitive.NewObjectID()
		teacher.TeacherID = teacher.ID.Hex()

		token, refreshToken, _ := helper.GenerateAllTokens(teacher.TeacherID, teacher.Name, teacher.Email, teacher.SchoolCode)
		teacher.Token = token
		teacher.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := teacherCollection.InsertOne(ctx, teacher)

		if insertErr != nil {
			msg := "teacher item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func LoginTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var teacher models.Teacher
		var foundteacher models.Teacher

		if err := c.BindJSON(&teacher); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := teacherCollection.FindOne(ctx, bson.M{"email": teacher.Email}).Decode(&foundteacher)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "teacher not found, login seems to be incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(teacher.Password, foundteacher.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(foundteacher.TeacherID, foundteacher.Name, foundteacher.Email, foundteacher.SchoolCode)

		helper.UpdateAllTeacherTokens(token, refreshToken, foundteacher.TeacherID)

		c.JSON(http.StatusOK, foundteacher)
	}
}

func LogOutTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		teacher_id := c.Param("teacher_id")

		_, err := teacherCollection.UpdateOne(ctx, bson.M{"teacher_id": teacher_id}, bson.M{"$set": bson.M{"tokens": ""}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while logging out"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}
