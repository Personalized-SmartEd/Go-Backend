package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Teacher struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" validate:"required"`
	Age          int                `bson:"age" validate:"required,gte=20"`
	Email        string             `bson:"email" validate:"required,email"`
	Password     string             `bson:"password" validate:"required,min=8"`
	TeacherID    string             `bson:"teacher_id" validate:"required"`
	Image        string             `bson:"image"`
	SchoolName   string             `bson:"school_name" validate:"required"`
	SchoolCode   string             `bson:"school_code" validate:"required"`
	Token        string             `bson:"token"`
	RefreshToken *string            `bson:"refresh_token"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
}
