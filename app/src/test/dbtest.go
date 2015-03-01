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
		ID           bson.ObjectId   `bson:"_id,omitempty"`
		Name         string          `bson:"name"`
		DateCreated  time.Time       `bson:"datecreated"`
		DateModified time.Time       `bson:"datemodigied"`
		Priority     int             `bson:"priority"`
		Buckets      []bson.ObjectId `bson:buckets"`
		Completed    bool            `bson:"completed"`
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

func insertBucket(bucket Bucket, collection *mgo.Collection) {
	err := collection.Insert(bucket)
	if err != nil {
		log.Fatal(err)
	}
}

func removeBucket(bucketName string, collection *mgo.Collection) {
	err := collection.Remove(bson.M{"name": bucketName})
	if err != nil {
		log.Fatal("error", err)
	}
}

/* Wrapper for calling other session tests */
func sessionTest(session *mgo.Session, f sessionFunc) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()

	f(sessionCopy)
}

/* Test inserting a bucket into the db */
func bucketTest(session *mgo.Session) {
	fmt.Printf("retrieving collection\n")
	collection := session.DB(TestDatabase).C("buckets")

	fmt.Printf("inserting into collection\n")
	bucket := Bucket{ID: bson.NewObjectId(), Name: "weekly"}
	insertBucket(bucket, collection)

	fmt.Printf("retrieving from collection\n")
	result := Bucket{}
	err := collection.Find(bson.M{"name": "weekly"}).One(&result)
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Println("Bucket:", result.Name)

	removeBucket("weekly", collection)
}

/* Test inserting a task into the db */
func taskTest(session *mgo.Session) {
	//bucketCollection := session.DB(TestDatabase).C("buckets")
	//taskCollection := session.DB(TestDatabase).C("tasks")

}

/* Main test suite */
func main() {
	mongoSession := dbSetup()
	sessionTest(mongoSession, bucketTest)
}
