package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AdminID  string             `bson:"admin_id" validate:"required"`
	Password string             `bson:"password" validate:"required,min=8"`
}
