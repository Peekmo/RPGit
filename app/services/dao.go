package services

import (
	"RPGithub/app/db"
	"RPGithub/app/model"
	"fmt"
	"github.com/revel/revel"
	"sort"
	"strings"
)

// Map reduce data (implements sort.Interface)
type MapReduceData []struct {
	Key   string `json:"key" bson:"_id"`
	Value int    `json:"value"`
}

// Len returns MapReduceData length
func (m MapReduceData) Len() int {
	return len(m)
}

// Swap swaps 2 values from MapReduceData
func (m MapReduceData) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Less checks if the first value is greater than the second
func (m MapReduceData) Less(i, j int) bool {
	return m[i].Value > m[j].Value
}

// GetUser gets a user from the database
func GetUser(username string) *model.User {
	var user *model.User

	userData := db.Database.Get(strings.ToLower(username), db.COLLECTION_USER)
	err := userData.One(&user)
	if err != nil {
		return nil
	}

	return user
}

// GetRepository gets a new repository from the database
func GetRepository(id int) *model.Repository {
	var repository *model.Repository

	repositoryData := db.Database.Get(id, db.COLLECTION_REPOSITORY)
	err := repositoryData.One(&repository)
	if err != nil {
		return nil
	}

	return repository
}

// GetUserRepositories gets a list of repositories from the given user
func GetUserRepositories(username string) []*model.Repository {
	var repositories []*model.Repository

	data := db.Database.GetQuery(map[string]string{"owner": strings.ToLower(username)}, db.COLLECTION_REPOSITORY)
	data.All(&repositories)

	return repositories
}

// UpdateUser updates the given user from the database
func UpdateUser(user *model.User) error {
	err := db.Database.Update(user.Id, user, db.COLLECTION_USER)
	if err != nil {
		revel.ERROR.Fatalf("Error while updating data : %s", err.Error())
		return err
	}

	return nil
}

// UpdateRepository updates the repository in the database
func UpdateRepository(repository *model.Repository) error {
	err := db.Database.Update(repository.Id, repository, db.COLLECTION_REPOSITORY)
	if err != nil {
		revel.ERROR.Fatalf("Error while updating data : %s", err.Error())
		return err
	}

	return nil
}

// RegisterRepository register in database a new repository
func RegisterRepository(repository *model.Repository) error {
	err := db.Database.Set(repository, db.COLLECTION_REPOSITORY)
	if err != nil {
		revel.ERROR.Fatalf("Error while saving new repository : %s", err.Error())
		return err
	}

	return nil
}

// RegisterRepository register in database a new repository
func RegisterUser(user *model.User) error {
	err := db.Database.Set(user, db.COLLECTION_USER)
	if err != nil {
		revel.ERROR.Fatalf("Error while saving new user : %s", err.Error())
		return err
	}

	return nil
}

// RegisterEventDay registers a new event
func RegisterEventDay(event *model.EventDay) error {
	err := db.Database.Set(event, db.COLLECTION_EVENT_DAY)
	if err != nil {
		revel.ERROR.Fatalf("Error while saving new event : %s", err.Error())
		return err
	}

	return nil
}

// RankingEventNumber gets from the daily events, the ranking by number of events
func RankingEventNumber(event string) (MapReduceData, error) {
	var result MapReduceData

	_, err := db.Database.MapReduce(
		fmt.Sprintf("function() { if (this.type == '%s') { emit(this.user, 1) } }", event),
		"function (key, values) { return Array.sum(values) }",
		db.COLLECTION_EVENT_DAY,
		&result,
	)

	fmt.Println(len(result))

	if err != nil {
		revel.ERROR.Fatalf("Error while mapreducing event number : %s", err.Error())
		return nil, err
	}

	sort.Sort(result)
	return result, nil
}

// ClearEventDay removes all elements from events collection
func ClearEventDay() error {
	_, err := db.Database.ClearCollection(db.COLLECTION_EVENT_DAY)
	if err != nil {
		revel.ERROR.Fatalf("Error while clearing events collection : %s", err.Error())
		return err
	}

	return nil
}
