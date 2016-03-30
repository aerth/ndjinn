package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aerth/ndjinn/components/database"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ObjectId   bson.ObjectId `bson:"_id"`
	MembershipLevel MembershipLevel    `db:"membership" bson:"membership"`
	NickName string        `db:"nickname" bson:"nickname"`
	Last_name  string        `db:"last_name" bson:"last_name"`
	Email      string        `db:"email" bson:"email"`
	Password   string        `db:"password" bson:"password"`
	Status_id  uint8         `db:"status_id" bson:"status_id"`
	Created_at time.Time     `db:"created_at" bson:"created_at"`
	Updated_at time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted    uint8         `db:"deleted" bson:"deleted"`

}

// User_status table contains every possible user status (active/inactive)
type User_status struct {
	Id         uint8     `db:"id" bson:"id"`
	Status     string    `db:"status" bson:"status"`
	Created_at time.Time `db:"created_at" bson:"created_at"`
	Updated_at time.Time `db:"updated_at" bson:"updated_at"`
	Deleted    uint8     `db:"deleted" bson:"deleted"`
	MembershipLevel uint8     `db:"membership" bson:"membership"`
}

var (
	ErrCode        = errors.New("Case statement in code is not correct.")
	ErrNoResult    = errors.New("Result not found.")
	ErrUnavailable = errors.New("Database is unavailable.")
)

type MembershipLevel int

const (
	NewMember MembershipLevel = iota
	PaidMember
	SuperMember
	Moderator
	Admin
)

// Id returns the user id
func (u *User) ID() string {
	r := ""

	switch database.ReadConfig().Type {

	case database.TypeBolt:
		r = u.ObjectId.Hex()
	}

	return r
}

// standardizeErrors returns the same error regardless of the database used
func standardizeError(err error) error {
	if err == sql.ErrNoRows || err == mgo.ErrNotFound {
		return ErrNoResult
	}

	return err
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	var err error

	result := User{}

	switch database.ReadConfig().Type {
	case database.TypeBolt:
		err = database.View("user", email, &result)
		if err != nil {
			err = ErrNoResult
		}
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

// UserCreate creates user
func UserCreate(nickname, last_name, email, password string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {

	case database.TypeBolt:
		user := &User{
			ObjectId:   bson.NewObjectId(),
			NickName: nickname,
			Last_name:  last_name,
			Email:      email,
			Password:   password,
			MembershipLevel: NewMember,
			Status_id:  1,
			Created_at: now,
			Updated_at: now,
			Deleted:    0,
		}

		err = database.Update("user", user.Email, &user)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// PromoteUser promotes a user
func UserPromote(user User) error {
	var err error

	now := time.Now()
	switch database.ReadConfig().Type {

	case database.TypeBolt:
		usernew := &User{

			MembershipLevel: PaidMember,
			Updated_at: now,
		}

		err = database.Update("user", user.Email, &usernew)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
