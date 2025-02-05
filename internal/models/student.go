package models

import (
	"time"

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
	Subjects        []string           `bson:"subjects" validate:"required"`
	Pace            string             `bson:"pace"`
	ClassNumber     int                `bson:"class_number"`
	ClassCode       string             `bson:"class_code"`
	Performance     float64            `bson:"performance"`
	PerformanceLvl  string             `bson:"performance_lvl"`
	PastPerformance []float64          `bson:"past_performance"`
	LearningStyle   string             `bson:"learning_style"`
	Token           string             `bson:"token"`
	RefreshToken    *string            `bson:"refresh_token"`
	Created_at      time.Time          `json:"created_at"`
	Updated_at      time.Time          `json:"updated_at"`
}
