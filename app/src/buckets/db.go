package buckets

import (
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
		DateModified time.Time       `bson:"datemodified"`
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
		log.Fatal(err)
	}
}

/* Query all buckets from the db */
func getBuckets(session *mgo.Session) {
	// Retrieve buckets collection.
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("buckets")

	var buckets []Bucket
	err := collection.Find(nil).All(&buckets)
	if err != nil {
		log.Printf("getBuckets ERROR: %s\n", err)
		return
	}

	log.Printf("bucketsQuery")
}

/* Create a bucket and insert it into the db */
func createBucket(session *mgo.Session, name string, tasks []string) *Bucket {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("buckets")

	tasksIds := make([]bson.ObjectId, len(tasks))
	for i, task := range tasks {
		tasksIds[i] = bson.ObjectIdHex(task)
	}
	bucket := Bucket{ID: bson.NewObjectId(), Name: name, Tasks: tasksIds}
	insertItem(bucket, collection)

	return &bucket
}

/* Retrieve a bucket from the db */
func getBucket(session *mgo.Session, id string) (*Bucket, error) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("buckets")

	bucket := Bucket{}
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&bucket)
	return &bucket, err
}

/* Remove a bucket from the db */
func removeBucket(session *mgo.Session, id string) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("buckets")

	removeItem(bson.M{"_id": bson.ObjectIdHex(id)}, collection)
}

/* Retrieve all tasks from the db */
func getTasks(session *mgo.Session) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("tasks")

	var tasks []Task
	err := collection.Find(nil).All(&tasks)
	if err != nil {
		log.Printf("getTasks ERROR: %s\n", err)
		return
	}
}

/* Insert a task into the db */
func createTask(session *mgo.Session, name string, priority int, buckets []string) *Task {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("tasks")

	bucketIds := make([]bson.ObjectId, len(buckets))
	for i, bucket := range buckets {
		bucketIds[i] = bson.ObjectIdHex(bucket)
	}

	task := Task{ID: bson.NewObjectId(), Name: name, DateCreated: time.Now(),
		DateModified: time.Now(), Priority: priority,
		Buckets: bucketIds, Completed: false}
	insertItem(task, collection)

	return &task
}

/* Retrieve a task from the db */
func getTask(session *mgo.Session, id string) (*Task, error) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("tasks")

	task := Task{}
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&task)
	return &task, err
}

/* Update a task in the db */
func updateTask(session *mgo.Session, id string, task *Task) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("tasks")

	err := collection.Update(bson.M{"_id": bson.ObjectIdHex(id)},
		task)
	if err != nil {
		log.Fatal("updateTask ERROR:", err)
	}
}

/* Remove a task from the db */
func removeTask(session *mgo.Session, id string) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("tasks")

	removeItem(bson.M{"_id": bson.ObjectIdHex(id)}, collection)
}
