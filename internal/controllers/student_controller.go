package controllers

import (
	"backend/internal/config"
	"backend/internal/helper"
	"backend/internal/models"
	"backend/internal/utils"
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing student items"})
			return
		}

		var allstudents []bson.M
		if err = result.All(ctx, &allstudents); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding student items"})
			return
		}

		c.JSON(http.StatusOK, allstudents)
	}
}

func GetStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		studentID := c.Param("student_id")
		var student models.Student

		err := studentCollection.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&student)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while retrieving student item"})
			return
		}
		c.JSON(http.StatusOK, student)
	}
}

func AddStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var student utils.StudentSignUp
		if err := c.BindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var inputStudent models.Student
		inputStudent.Name = student.Name
		inputStudent.Age = student.Age
		inputStudent.Password = student.Password
		inputStudent.Email = student.Email
		inputStudent.Image = student.Image
		inputStudent.SchoolName = student.SchoolName
		inputStudent.SchoolCode = student.SchoolCode
		inputStudent.Subjects = student.Subjects
		inputStudent.Class = student.Class

		inputStudent.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		inputStudent.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		inputStudent.ID = primitive.NewObjectID()
		inputStudent.StudentID = inputStudent.ID.Hex()

		inputStudent.Performance = 0.0
		inputStudent.PerformanceLvl = "beginner"
		inputStudent.PastPerformance = []float64{0.0}
		inputStudent.LearningStyle = ""
		inputStudent.Pace = "slow"

		validationErr := validate.Struct(inputStudent)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := studentCollection.CountDocuments(ctx, bson.M{"email": inputStudent.Email})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			return
		}

		inputStudent.Password = HashPassword(inputStudent.Password)

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		resultInsertionNumber, insertErr := studentCollection.InsertOne(ctx, inputStudent)

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

		var student utils.StudentLogin
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

		foundstudent.Token = token
		foundstudent.RefreshToken = &refreshToken

		helper.UpdateAllStudentTokens(token, refreshToken, foundstudent.StudentID)

		c.JSON(http.StatusOK, foundstudent)
	}
}

func UpdateStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		studentID := c.Param("student_id")

		var studentUpdate utils.StudentSignUp
		if err := c.BindJSON(&studentUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := bson.M{
			"name":        studentUpdate.Name,
			"age":         studentUpdate.Age,
			"password":    HashPassword(studentUpdate.Password),
			"email":       studentUpdate.Email,
			"image":       studentUpdate.Image,
			"school_name": studentUpdate.SchoolName,
			"school_code": studentUpdate.SchoolCode,
			"subjects":    studentUpdate.Subjects,
			"class":       studentUpdate.Class,
			"updated_at": func() time.Time {
				t, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
				return t
			}(),
		}

		result, err := studentCollection.UpdateOne(
			ctx,
			bson.M{"student_id": studentID},
			bson.M{"$set": update},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while updating the student item"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "student updated successfully", "result": result})
	}
}

func DeleteStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		studentID := c.Param("student_id")

		result, err := studentCollection.DeleteOne(ctx, bson.M{"student_id": studentID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while deleting the student"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "student deleted successfully"})
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
