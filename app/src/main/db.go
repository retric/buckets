package main

import (
	//"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const (
	MongoDBHosts = "localhost:27017"
	AuthDatabase = "test"
	TestDatabase = "test"
)

type (
	// Task
	Task struct {
		ID           bson.ObjectId   `bson:"_id"`
		Name         string          `bson:"name"`
		DateCreated  time.Time       `bson:"datecreated"`
		DateModified time.Time       `bson:"datemodigied"`
		Priority     int             `bson:"priority"`
		Buckets      []bson.ObjectId `bson:buckets"`
		Completed    bool            `bson:"completed"`
	}

	// Bucket
	Bucket struct {
		ID    bson.ObjectId   `bson:"_id"`
		Name  string          `bson:"name"`
		Tasks []bson.ObjectId `bson:"tasks"`
	}
)

func dbSetup() *mgo.Session {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
	}

	// Init session with DB
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	mongoSession.SetMode(mgo.Monotonic, true)

	return mongoSession
}

func bucketsQuery(session *mgo.Session) {

	// Request socket connection from session.
	// Close session when function is done and return connection to the pool.
	sessionCopy := session.Copy()
	defer sessionCopy.Close()

	// Retrieve buckets collection.
	collection := sessionCopy.DB(TestDatabase).C("buckets")

	var buckets []Bucket
	err := collection.Find(nil).All(&buckets)
	if err != nil {
		log.Printf("bucketsQuery ERROR: %s\n", err)
		return
	}

	log.Printf("bucketsQuery")
}