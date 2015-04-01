package buckets

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
	"time"
)

const (
	TestDatabase = "test"
)

type (
	sessionFunc func(*mgo.Session)
)

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

/* Test inserting a task into the db */
func taskTest(session *mgo.Session) {
	fmt.Printf("taskTest: retrieving collection\n")
	bucketCollection := session.DB(TestDatabase).C("buckets")
	taskCollection := session.DB(TestDatabase).C("tasks")

	bucket := Bucket{ID: bson.NewObjectId(), Name: "weekly"}
	insertItem(bucket, bucketCollection)

	task := Task{ID: bson.NewObjectId(), Name: "read", Priority: 1,
		DateCreated: time.Now().Local(), DateModified: time.Now().Local(),
		Buckets: []bson.ObjectId{}, Completed: false}
	task.Buckets = append(task.Buckets, bucket.ID)
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

/* Test suite for Buckets */
func TestBucket(t *testing.T) {
	mongoSession := dbSetup()
	bucket := CreateBucketTest(mongoSession)
	id := bucket.ID.Hex()

	GetBucketTest(mongoSession, id)
	RemoveBucketTest(mongoSession, id)
}

/* Test retrieving a bucket */
func GetBucketTest(session *mgo.Session, id string) {
	_, err := getBucket(session, id)
	if err != nil {
		log.Fatal("error: bucket not found")
	}
}

/* Test creating a bucket */
func CreateBucketTest(session *mgo.Session) *Bucket {
	bucket := createBucket(session, "weekly", []string{"54f41e6a5786752068000003"})
	return bucket
}

/* Test removing a bucket */
func RemoveBucketTest(session *mgo.Session, id string) {
	removeBucket(session, id)
	_, err := getBucket(session, id)
	if err == nil {
		log.Fatal("error: bucket found after remove attempt")
	}
}

/* Test suite for Tasks */
func TestTask(t *testing.T) {
	mongoSession := dbSetup()
	task := CreateTaskTest(mongoSession)
	id := task.ID.Hex()

	GetTaskTest(mongoSession, id)
	UpdateTaskTest(mongoSession, id, task)
	RemoveTaskTest(mongoSession, id)
}

/* Test for creating a task */
func CreateTaskTest(session *mgo.Session) *Task {
	task := createTask(session, "running", 1, []string{"55145cdb5786751845000001"})
	return task
}

/* Test for retrieving a task */
func GetTaskTest(session *mgo.Session, id string) {
	_, err := getTask(session, id)
	if err != nil {
		log.Fatal("error: task not found")
	}
}

/* Test for updating a task */
func UpdateTaskTest(session *mgo.Session, id string, task *Task) {
	priority := task.Priority
	task.Priority = priority + 1
	updateTask(session, id, task)
	task2, err := getTask(session, id)
	if err != nil {
		log.Fatal("error: task not found after update")
	} else if task2.Priority != priority+1 {
		log.Fatal("error: priority not updated")
	}
}

/* Test for removing a task */
func RemoveTaskTest(session *mgo.Session, id string) {
	removeTask(session, id)
	_, err := getTask(session, id)
	if err == nil {
		log.Fatal("error: task found after remove attempt")
	}
}
