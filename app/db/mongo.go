package db

import (
	"fmt"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
)

type Mongo struct {
	db *mgo.Database
}

var (
	Database *Mongo
)

// Collections available
const (
	COLLECTION_USER         = "User"
	COLLECTION_EVENT_DAY    = "EventDay"
	COLLECTION_REPOSITORY   = "Repository"
	COLLECTION_ORGANIZATION = "Organization"
)

// Set a document in the database
func (this *Mongo) Set(doc interface{}, collection string) error {
	return this.db.C(collection).Insert(doc)
}

// Get a document by its identifier
func (this *Mongo) Get(id interface{}, collection string) *mgo.Query {
	return this.db.C(collection).FindId(id)
}

// Update the given document
func (this *Mongo) Update(key, value interface{}, collection string) error {
	return this.db.C(collection).UpdateId(key, value)
}

// GetQuery gets all documents following the givne query
func (this *Mongo) GetQuery(query interface{}, collection string) *mgo.Query {
	return this.db.C(collection).Find(query)
}

// UpdateQuery updates documents with the given query
func (this *Mongo) UpdateQuery(query, data interface{}, collection string) error {
	return this.db.C(collection).Update(query, data)
}

// ClearCollection removes all documents from the given collection
func (this *Mongo) ClearCollection(collection string) (*mgo.ChangeInfo, error) {
	return this.db.C(collection).RemoveAll(map[string]string{})
}

// MapReduce executes the given map reduce function
func (this *Mongo) MapReduce(mapfunc, reduce, collection string, result interface{}) (*mgo.MapReduceInfo, error) {
	job := &mgo.MapReduce{
		Map:    mapfunc,
		Reduce: reduce,
	}

	return this.db.C(collection).Find(nil).MapReduce(job, result)
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

	db := session.DB(revel.Config.StringDefault("mongo.database", "RPGithub"))

	Database = &Mongo{db}
}
