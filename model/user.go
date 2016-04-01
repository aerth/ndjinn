package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aerth/ndjinn/components/database"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MembershipLevel int

const (
	NewMember MembershipLevel = iota
	PaidMember
	SuperMember
	Moderator
	Admin
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ObjectId        bson.ObjectId   `bson:"_id"`
	Id              uint32          `db:"id" bson:"id,omitempty"` // Don't use Id, use ID() instead for consistency with MongoDB
	NickName        string          `db:"NickName" bson:"NickName"`
	MembershipLevel MembershipLevel `db:"membershiplevel" bson:"membershiplevel"`
	Email           string          `db:"email" bson:"email"`
	Password        string          `db:"password" bson:"password"`
	Status_id       uint8           `db:"status_id" bson:"status_id"`
	Created_at      time.Time       `db:"created_at" bson:"created_at"`
	Updated_at      time.Time       `db:"updated_at" bson:"updated_at"`
	Deleted         uint8           `db:"deleted" bson:"deleted"`
}

// User_status table contains every possible user status (active/inactive)
type User_status struct {
	Id         uint8     `db:"id" bson:"id"`
	Status     string    `db:"status" bson:"status"`
	Created_at time.Time `db:"created_at" bson:"created_at"`
	Updated_at time.Time `db:"updated_at" bson:"updated_at"`
	Deleted    uint8     `db:"deleted" bson:"deleted"`
}

var (
	ErrCode        = errors.New("Case statement in code is not correct.")
	ErrNoResult    = errors.New("Result not found.")
	ErrUnavailable = errors.New("Database is unavailable.")
)

// Id returns the user id
func (u *User) ID() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		r = fmt.Sprintf("%v", u.Id)
	case database.TypeMongoDB:
		r = u.ObjectId.Hex()
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
	case database.TypeMySQL:
		err = database.Sql.Get(&result, "SELECT id, password, status_id, NickName FROM user WHERE email = ? LIMIT 1", email)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
			err = c.Find(bson.M{"email": email}).One(&result)
		} else {
			err = ErrUnavailable
		}
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
func UserCreate(NickName, MembershipLevel, email, password string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.Sql.Exec("INSERT INTO user (NickName, MembershipLevel, email, password) VALUES (?,?,?,?)", NickName,
			MembershipLevel, email, password)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("user")

			user := &User{
				ObjectId:        bson.NewObjectId(),
				NickName:        NickName,
				MembershipLevel: NewMember,
				Email:           email,
				Password:        password,
				Status_id:       1,
				Created_at:      now,
				Updated_at:      now,
				Deleted:         0,
			}
			err = c.Insert(user)
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		user := &User{
			ObjectId:        bson.NewObjectId(),
			NickName:        NickName,
			MembershipLevel: PaidMember,
			Email:           email,
			Password:        password,
			Status_id:       1,
			Created_at:      now,
			Updated_at:      now,
			Deleted:         0,
		}

		err = database.Update("user", user.Email, &user)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

// PromoteUser promotes a user
func UserPromote(user User) (err error) {

	now := time.Now()
	switch database.ReadConfig().Type {

	case database.TypeBolt:
		usernew := &User{
			ObjectId:        bson.ObjectIdHex(user.ID()),
			NickName:        user.NickName,
			MembershipLevel: PaidMember,
			Email:           user.Email,
			Password:        user.Password,
			Status_id:       1,
			Created_at:      user.Created_at,
			Updated_at:      now,
			Deleted:         user.Deleted,
		}

		err = database.Update("user", user.Email, &usernew)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
