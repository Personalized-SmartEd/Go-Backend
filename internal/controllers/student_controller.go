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
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var studentCollection *mongo.Collection = config.OpenCollection(config.Client, "student")
var validate = validator.New()

func GetStudents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := studentCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing student items"})
			return
		}

		var allstudents []bson.M
		if err = result.All(ctx, &allstudents); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while decoding food items"})
			return
		}

		c.JSON(http.StatusOK, allstudents)
	}
}

func GetStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		student_id := c.Param("student_id")

		var student models.Student

		err := studentCollection.FindOne(ctx, bson.M{"student_id": student_id}).Decode(&student)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing student items"})
			return
		}
		c.JSON(http.StatusOK, student)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var student models.Student

		if err := c.BindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(student)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := studentCollection.CountDocuments(ctx, bson.M{"email": student.Email})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := HashPassword(student.Password)
		student.Password = password

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exsits"})
			return
		}

		student.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		student.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		student.ID = primitive.NewObjectID()
		student.StudentID = student.ID.Hex()

		token, refreshToken, _ := helper.GenerateAllTokens(student.StudentID, student.Name, student.Email, student.Class)
		student.Token = token
		student.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := studentCollection.InsertOne(ctx, student)

		if insertErr != nil {
			msg := "student item was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var student models.Student
		var foundstudent models.Student

		if err := c.BindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := studentCollection.FindOne(ctx, bson.M{"email": student.Email}).Decode(&foundstudent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "student not found, login seems to be incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(student.Password, foundstudent.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(foundstudent.StudentID, foundstudent.Name, foundstudent.Email, foundstudent.Class)

		helper.UpdateAllTokens(token, refreshToken, foundstudent.StudentID)

		c.JSON(http.StatusOK, foundstudent)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(studentPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(studentPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "login or password is incorrect"
		check = false
	}
	return check, msg
}
