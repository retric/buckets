package buckets

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

/* Query all buckets from the db */
func getBuckets(session *mgo.Session) {
	// Request socket connection from session.
	// Close session when function is done and return connection to the pool.
	sessionCopy := session.Copy()
	defer sessionCopy.Close()

	// Retrieve buckets collection.
	collection := sessionCopy.DB(AuthDatabase).C("buckets")

	var buckets []Bucket
	err := collection.Find(nil).All(&buckets)
	if err != nil {
		log.Printf("bucketsQuery ERROR: %s\n", err)
		return
	}

	log.Printf("bucketsQuery")
}

/* Create a bucket and insert it into the db */
func createBucket(session *mgo.Session, name string, tasks []string) *Bucket {
	collection := session.DB(AuthDatabase).C("buckets")

	tasksIds := make([]bson.ObjectId, len(tasks))
	for i, task := range tasks {
		tasksIds[i] = bson.ObjectIdHex(task)
	}
	bucket := Bucket{ID: bson.NewObjectId(), Name: name, Tasks: tasksIds}
	insertItem(bucket, collection)

	return &bucket
}

/* Retrieve a bucket from the db */
func getBucket(session *mgo.Session, id string) *Bucket {
	collection := session.DB(AuthDatabase).C("buckets")

	bucket := Bucket{}
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&bucket)
	if err != nil {
		log.Fatal("getBucket ERROR:", err)
	}
	return &bucket
}

/* Remove a bucket from the db */
func removeBucket(session *mgo.Session, id string) {
	collection := session.DB(AuthDatabase).C("buckets")

	err := collection.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		log.Fatal("removeBucket ERROR:", err)
	}
}

func getTasks(session *mgo.Session) {

}

func getTask(session *mgo.Session, id string) {

}
