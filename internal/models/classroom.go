package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Classroom struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TeacherID   string             `bson:"teacher_id" validate:"required"`
	SchoolCode  string             `bson:"school_code" validate:"required"`
	Students    []string           `bson:"students"`
	ClassNumber string             `bson:"class_number" validate:"required"`
	ClassCode   string             `bson:"class_code" validate:"required"`
}
