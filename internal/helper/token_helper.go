package helper

import (
	"backend/internal/config"
	"backend/internal/utils"
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

var studentCollection *mongo.Collection = config.OpenCollection(config.Client, "student")
var teacherCollection *mongo.Collection = config.OpenCollection(config.Client, "teacher")

func GenerateAllTokens(studentID string, name string, email string, class string) (signedToken string, signedRefreshToken string, err error) {
	claims := &utils.SignedDetailsStudent{
		StudentID: studentID,
		Name:      name,
		Email:     email,
		Class:     class,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1000)).Unix(),
		},
	}

	refreshClaims := &utils.SignedDetailsStudent{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(400)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAllStudentTokens(signedToken string, signedRefreshToken string, studentId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.M{"student_id": studentId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := studentCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

}

func UpdateAllTeacherTokens(signedToken string, signedRefreshToken string, teacherId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updatedAt := time.Now()
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updatedAt})

	upsert := true
	filter := bson.M{"teacher_id": teacherId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := teacherCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	if err != nil {
		log.Panic(err)
		return
	}
}

func ValidateStudentToken(signedToken string) (claims *utils.SignedDetailsStudent, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&utils.SignedDetailsStudent{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	claims, ok := token.Claims.(*utils.SignedDetailsStudent)
	if !ok {
		msg = "the token is invalid"
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		msg = err.Error()
		return
	}

	return claims, msg

}

func ValidateTeacherToken(signedToken string) (claims *utils.SignedDetailsTeacher, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&utils.SignedDetailsTeacher{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	claims, ok := token.Claims.(*utils.SignedDetailsTeacher)
	if !ok {
		msg = "the token is invalid"
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		msg = err.Error()
		return
	}

	return claims, msg

}
