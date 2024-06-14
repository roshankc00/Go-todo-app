package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct{
	ID				primitive.ObjectID		`bson:"_id"`
	Title		*string					`json:"title" validate:"required,min=2,max=100"`
	Description		*string					`json:"description" validate:"required,min=2,max=100"`
	Status		*string					`json:"status" validate:"eq=DONE|eq=PENDING"`
	Created_at		time.Time				`json:"created_at"`
	Updated_at		time.Time				`json:"updated_at"`
	User_uid		string					`json:"user_id"`
}