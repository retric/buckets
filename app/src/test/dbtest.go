package test

import (
	"fmt"
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

/* Document structs */
type (
	// Task
	Task struct {
		ID           bson.ObjectId `bson:"_id,omitempty"`
		Name         string        `bson:"name"`
		DateCreated  time.Time     `bson:"datecreated"`
		DateModified time.Time     `bson:"datemodified"`
		Priority     int           `bson:"priority"`
		Buckets      []string      `bson:buckets"`
		Completed    bool          `bson:"completed"`
	}

	// Buckets
	Bucket struct {
		ID    bson.ObjectId   `bson:"_id,omitempty"`
		Name  string          `bson:"name"`
		Tasks []bson.ObjectId `bson:"tasks"`
	}

	sessionFunc func(*mgo.Session)
)

/* Initialize session with database */
func dbSetup() *mgo.Session {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	mongoSession.SetMode(mgo.Monotonic, true)

	return mongoSession
}

/* Query all buckets from the db */
func bucketsQuery(session *mgo.Session) {
	// Retrieve buckets collection.
	collection := session.DB(TestDatabase).C("buckets")

	var buckets []Bucket
	err := collection.Find(nil).All(&buckets)
	if err != nil {
		log.Printf("bucketsQuery ERROR: %s\n", err)
		return
	}

	log.Printf("bucketsQuery")
}

func insertItem(doc interface{}, collection *mgo.Collection) {
	err := collection.Insert(doc)
	if err != nil {
		log.Fatal(err)
	}
}

func removeItem(doc interface{}, collection *mgo.Collection) {
	err := collection.Remove(doc)
	if err != nil {
		fmt.Println("error", err)
	}
}

/* Wrapper for calling other session tests */
func sessionWrap(session *mgo.Session, f sessionFunc) {
	// Request socket connection from session.
	// Close session when function is done and return connection to the pool.
	sessionCopy := session.Copy()
	defer sessionCopy.Close()

	f(sessionCopy)
}

/* Test inserting a bucket into the db */
func bucketTest(session *mgo.Session) {
	fmt.Printf("bucketTest: retrieving collection\n")
	collection := session.DB(TestDatabase).C("buckets")

	fmt.Printf("inserting bucket into collection\n")
	bucket := Bucket{ID: bson.NewObjectId(), Name: "weekly"}
	insertItem(bucket, collection)

	fmt.Printf("retrieving bucket from collection\n")
	result := Bucket{}
	err := collection.Find(bson.M{"name": "weekly"}).One(&result)
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Println("Bucket:", result.Name)
	fmt.Printf("removing bucket from collection\n")
	removeItem(bson.M{"name": "weekly"}, collection)
}

/* Test inserting a task into the db */
func taskTest(session *mgo.Session) {
	fmt.Printf("taskTest: retrieving collection\n")
	bucketCollection := session.DB(TestDatabase).C("buckets")
	taskCollection := session.DB(TestDatabase).C("tasks")

	bucket := Bucket{ID: bson.NewObjectId(), Name: "weekly"}
	insertItem(bucket, bucketCollection)

	task := Task{ID: bson.NewObjectId(), Name: "read", Priority: 1,
		DateCreated: time.Now().Local(), DateModified: time.Now().Local(),
		Buckets: []string{}, Completed: false}
	task.Buckets = append(task.Buckets, bucket.Name)
	fmt.Printf("taskTest: inserting task into collection")
	insertItem(task, taskCollection)

	fmt.Printf("retrieving task from collection")
	result := Task{}
	err := taskCollection.Find(bson.M{"name": "read"}).One(&result)
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Println("Task:", result.Name)

	removeItem(bson.M{"name": "weekly"}, bucketCollection)
	removeItem(task, taskCollection)
}

/* Main test suite */
func TestMain() {
	mongoSession := dbSetup()
	sessionWrap(mongoSession, bucketTest)
	sessionWrap(mongoSession, taskTest)
}
