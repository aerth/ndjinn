package model

import (
	"fmt"
	"time"

	"github.com/aerth/ndjinn/components/database"

	"gopkg.in/mgo.v2/bson"
)

// Listing is separate from Users. Cool!
type Listing struct {
	ObjectId   bson.ObjectId `bson:"_id"`
	Id         uint32        `db:"id" bson:"id,omitempty"` // Don't use Id, use ID() instead for consistency with MongoDB
	Email      string        `db:"email" bson:"email"`
	Phone      string        `db:"phone" bson:"phone"`
	Content    string        `db:"content" bson:"content"`
	Deleted    uint8         `db:"deleted" bson:"deleted"`
	Created_at time.Time     `db:"created_at" bson:"created_at"`
	Updated_at time.Time     `db:"updated_at" bson:"updated_at"`
}

// Return Listing ID
func (t *Listing) ID() string {
	r := ""

	switch database.ReadConfig().Type {

	case database.TypeBolt:
		r = t.ObjectId.Hex()
	}

	return r
}

// ListingCreate creates listing
func ListingCreate(email, phone, content string) error {
	var err error

	now := time.Now()

	listing := &Listing{
		ObjectId:   bson.NewObjectId(),
		Email:      email,
		Phone:      phone,
		Content:    content,
		Created_at: now,
		Updated_at: now,
		Deleted:    0,
	}

	err = database.Update("listing", "listing", &listing)
	err = database.Update("listing", "listing", &listing)
	fmt.Println(database.View("listing", "listing", &listing))
	return standardizeError(err)
}
