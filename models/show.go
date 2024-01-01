package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Show struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	ReleaseDate time.Time          `json:"releaseDate,omitempty" bson:"release_date,omitempty"`
	Description string			   `json:"description,omitempty" bson:"description,omitempty"`
	Performers  []string           `json:"performers,omitempty" bson:"performers,omitempty"`
	Accepted    bool               `json:"accepted,omitempty" bson:"accepted,omitempty"`
}