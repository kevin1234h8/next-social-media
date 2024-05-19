package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Username  string             `bson:"username" json:"username"  binding:"required,min=3,max=20"`
	Email     string             `bson:"email" json:"email" binding:"email,required"`
	Password  string             `bson:"password" json:"password" binding:"required"`
	IsAdmin   bool               `bson:"is_admin" json:"is_admin"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
