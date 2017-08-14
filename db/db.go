package db

import (
	"os"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

// MONGODB_URI the variable where address of database is stored
var (
	MONGODB_URI = os.Getenv("MONGODB_URI")
)

// DB struct has all the methods for performing operations on database
type DB struct {
	Sess *mgo.Session
}

// User struct is the model that is stored in database
type User struct {
	ChatID   int64 `bson:"_id"`
	Username string
}

// Init initialises the package
func Init() (*DB, error) {

	sess, err := mgo.Dial(MONGODB_URI)
	if err != nil {
		return nil, err
	}

	db := &DB{Sess: sess}

	return db, nil
}

// AddSubscriber adds a information about  subscriber to the database
func (d *DB) AddSubscriber(u *User) error {
	s := d.Sess.Copy()

	c := s.DB("ch").C("subs")

	err := c.Insert(u)

	if err != nil {
		return err
	}

	return nil
}

// RemoveSubscriber removes the information of said subscriber from the database
func (d *DB) RemoveSubscriber(u *User) error {
	s := d.Sess.Copy()

	c := s.DB("ch").C("subs")

	err := c.Remove(u)
	if err != nil {
		return err
	}
	return nil
}

// FetchSubscribers returns a slice of all the subscribers in Databaase
func (d *DB) FetchSubscribers() ([]User, error) {
	s := d.Sess.Copy()
	c := s.DB("ch").C("subs")

	var results []User

	err := c.Find(bson.M{}).All(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
