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

	BucketPart struct {
		Name  string
		Tasks []string
	}

	TaskPart struct {
		Name         string
		DateCreated  time.Time
		DateModified time.Time
		Priority     int
		Buckets      []string
		Completed    bool
	}
)

/* Initialize session with database */
func DbSetup() *mgo.Session {
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
func getBuckets(session *mgo.Session) []Bucket {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("buckets")

	var buckets []Bucket
	err := collection.Find(nil).All(&buckets)
	if err != nil {
		log.Fatal("getBuckets ERROR: %s\n", err)
		return nil
	}
	return buckets
}

/* Create a bucket and insert it into the db */
func createBucket(session *mgo.Session, bForm BucketPart) *Bucket {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("buckets")

	tasksIds := make([]bson.ObjectId, len(bForm.Tasks))
	for i, task := range bForm.Tasks {
		tasksIds[i] = bson.ObjectIdHex(task)
	}
	bucket := Bucket{ID: bson.NewObjectId(), Name: bForm.Name, Tasks: tasksIds}
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

/* Update a bucket in the db */
func updateBucket(session *mgo.Session, id string, bForm BucketPart) (*Bucket, error) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("buckets")

	tasksIds := make([]bson.ObjectId, len(bForm.Tasks))
	for i, task := range bForm.Tasks {
		tasksIds[i] = bson.ObjectIdHex(task)
	}
	bucketUpdate := bson.M{"name": bForm.Name, "tasks": tasksIds}
	err := collection.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bucketUpdate})
	if err != nil {
		log.Fatal("updateBucket ERROR:", err)
	}
	bucket := Bucket{}
	err = collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&bucket)
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
func createTask(session *mgo.Session, tForm TaskPart) *Task {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("tasks")

	bucketIds := make([]bson.ObjectId, len(tForm.Buckets))
	for i, bucket := range tForm.Buckets {
		bucketIds[i] = bson.ObjectIdHex(bucket)
	}

	task := Task{ID: bson.NewObjectId(), Name: tForm.Name, DateCreated: time.Now(),
		DateModified: time.Now(), Priority: tForm.Priority,
		Buckets: bucketIds, Completed: tForm.Completed}
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
func updateTask(session *mgo.Session, id string, tForm TaskPart) (*Task, error) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("tasks")

	bucketIds := make([]bson.ObjectId, len(tForm.Buckets))
	for i, bucket := range tForm.Buckets {
		bucketIds[i] = bson.ObjectIdHex(bucket)
	}

	taskUpdate := bson.M{"name": tForm.Name,
		"priority":     tForm.Priority,
		"datemodified": time.Now(),
		"buckets":      bucketIds,
		"completed":    tForm.Completed}
	err := collection.Update(bson.M{"_id": bson.ObjectIdHex(id)},
		bson.M{"$set": taskUpdate})
	if err != nil {
		log.Fatal("updateTask ERROR:", err)
	}
	task := Task{}
	err = collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&task)
	return &task, err
}

/* Remove a task from the db */
func removeTask(session *mgo.Session, id string) {
	sessionCopy := session.Copy()
	defer sessionCopy.Close()
	collection := sessionCopy.DB(AuthDatabase).C("tasks")

	removeItem(bson.M{"_id": bson.ObjectIdHex(id)}, collection)
}
