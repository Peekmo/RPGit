package db

import (
	"fmt"
	"time"

	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
)

type Mongo struct {
	Db      *mgo.Database
	Session *mgo.Session
}

var (
	Database *Mongo
)

// Collections available
const (
	COLLECTION_USER         = "User"
	COLLECTION_BLACKLIST    = "Blacklist"
	COLLECTION_EVENT_DAY    = "EventDay"
	COLLECTION_REPOSITORY   = "Repository"
	COLLECTION_ORGANIZATION = "Organization"
)

// Set a document in the database
func (this *Mongo) Set(doc interface{}, collection string) error {
	return this.Db.C(collection).Insert(doc)
}

// Get a document by its identifier
func (this *Mongo) Get(id interface{}, collection string) *mgo.Query {
	return this.Db.C(collection).FindId(id)
}

// Update the given document
func (this *Mongo) Update(key, value interface{}, collection string) error {
	return this.Db.C(collection).UpdateId(key, value)
}

// GetQuery gets all documents following the givne query
func (this *Mongo) GetQuery(query interface{}, collection string) *mgo.Query {
	return this.Db.C(collection).Find(query)
}

// UpdateQuery updates documents with the given query
func (this *Mongo) UpdateQuery(query, data interface{}, collection string) error {
	return this.Db.C(collection).Update(query, data)
}

// ClearCollection removes all documents from the given collection
func (this *Mongo) ClearCollection(collection string) (*mgo.ChangeInfo, error) {
	return this.Db.C(collection).RemoveAll(map[string]string{})
}

// Remove removes elements from the given collection following the query
func (this *Mongo) Remove(query interface{}, collection string) (*mgo.ChangeInfo, error) {
	return this.Db.C(collection).RemoveAll(query)
}

// MapReduce executes the given map reduce function
func (this *Mongo) MapReduce(mapfunc, reduce, collection string, result interface{}) (*mgo.MapReduceInfo, error) {
	job := &mgo.MapReduce{
		Map:    mapfunc,
		Reduce: reduce,
	}

	return this.Db.C(collection).Find(nil).MapReduce(job, result)
}

// Changes db's session (timeout reason)
func (this *Mongo) InitSession() {
	this.Db.Session = this.Session.Copy()
}

// InitDatabse initialize the mongodb session
func InitDatabase() {
	var url string
	address := revel.Config.StringDefault("mongo.address", "127.0.0.1")
	port := revel.Config.StringDefault("mongo.port", "27017")
	url = fmt.Sprintf("mongodb://%s:%s", address, port)

	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	session.SetSocketTimeout(1 * time.Hour)
	db := session.DB(revel.Config.StringDefault("mongo.database", "RPGithub"))

	Database = &Mongo{db, db.Session}
}
