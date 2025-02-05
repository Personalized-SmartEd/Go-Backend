package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" validate:"required"`
	Age             int                `bson:"age" validate:"required,gte=5"`
	Password        string             `bson:"password" validate:"required,min=8"`
	Email           string             `bson:"email" validate:"required,email"`
	StudentID       string             `bson:"student_id" validate:"required"`
	Image           string             `bson:"image"`
	SchoolName      string             `bson:"school_name" validate:"required"`
	SchoolCode      string             `bson:"school_code" validate:"required"`
	Subjects        []string           `bson:"subjects" validate:"required,subjects"`
	Pace            string             `bson:"pace" validate:"oneof=slow moderate fast"`
	ClassNumber     int                `bson:"class_number"`
	ClassCode       string             `bson:"class_code"`
	Performance     float64            `bson:"performance"`
	PerformanceLvl  string             `bson:"performance_lvl" validate:"oneof=beginner intermediate advanced"`
	PastPerformance []float64          `bson:"past_performance"`
	LearningStyle   string             `bson:"learning_style" validate:"oneof=visual auditory reading_writing kinesthetic"`
	Token           string             `bson:"token"`
	RefreshToken    *string            `bson:"refresh_token"`
	Created_at      time.Time          `json:"created_at"`
	Updated_at      time.Time          `json:"updated_at"`
}

var validSubjects = []string{"math", "science", "english", "social-science", "hindi"}

func validateSubjects(fl validator.FieldLevel) bool {
	subjects, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}
	for _, subject := range subjects {
		if !contains(validSubjects, subject) {
			return false
		}
	}
	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
