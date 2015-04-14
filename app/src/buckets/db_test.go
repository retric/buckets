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
	mongoSession := DbSetup()
	bucket := CreateBucketTest(mongoSession)
	id := bucket.ID.Hex()

	GetBucketTest(mongoSession, id)
	UpdateBucketTest(mongoSession, id)
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
	bucketPart := BucketPart{Name: "weekly", Tasks: []string{"54f41e6a5786752068000003"}}
	bucket := createBucket(session, bucketPart)
	return bucket
}

/* Test updating a bucket */
func UpdateBucketTest(session *mgo.Session, id string) {
	bucketPart := BucketPart{Name: "monthly", Tasks: []string{"54f41e6a5786752068000004"}}
	updateBucket(session, id, bucketPart)
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
	mongoSession := DbSetup()
	task := CreateTaskTest(mongoSession)
	id := task.ID.Hex()

	GetTaskTest(mongoSession, id)
	UpdateTaskTest(mongoSession, id, task)
	RemoveTaskTest(mongoSession, id)
}

/* Test for creating a task */
func CreateTaskTest(session *mgo.Session) *Task {
	taskPart := TaskPart{Name: "running", Priority: 1, Buckets: []string{"55145cdb5786751845000001"}}
	task := createTask(session, taskPart)
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
	tForm := TaskPart{Name: task.Name,
		DateCreated:  task.DateCreated,
		DateModified: task.DateModified,
		Priority:     task.Priority + 1,
		Buckets:      []string{"54f41e6a5786752068000003"},
		Completed:    task.Completed}
	updateTask(session, id, tForm)
	task2, err := getTask(session, id)
	if err != nil {
		log.Fatal("error: task not found after update")
	} else if task2.Priority != task.Priority+1 {
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
