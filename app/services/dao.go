package services

import (
	"RPGithub/app/db"
	"RPGithub/app/model"
	"fmt"
	"sort"
	"strings"

	"github.com/revel/revel"
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

// InitDatabase starts the database
func InitDatabase() {
	db.InitDatabase()
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

// IsBlacklisted checks if the given user is blacklisted or not
func IsBlacklisted(name string) bool {
	var blacklist *model.Blacklist

	blacklistData := db.Database.Get(strings.ToLower(name), db.COLLECTION_BLACKLIST)
	err := blacklistData.One(&blacklist)
	if err != nil {
		return false
	}

	return true
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
		revel.ERROR.Printf("Error while updating data : %s", err.Error())
		return err
	}

	return nil
}

// UpdateRepository updates the repository in the database
func UpdateRepository(repository *model.Repository) error {
	err := db.Database.Update(repository.Id, repository, db.COLLECTION_REPOSITORY)
	if err != nil {
		revel.ERROR.Printf("Error while updating data : %s", err.Error())
		return err
	}

	return nil
}

// RegisterRepository register in database a new repository
func RegisterRepository(repository *model.Repository) error {
	err := db.Database.Set(repository, db.COLLECTION_REPOSITORY)
	if err != nil {
		revel.ERROR.Printf("Error while saving new repository : %s", err.Error())
		return err
	}

	return nil
}

// RegisterRepository register in database a new repository
func RegisterUser(user *model.User) error {
	err := db.Database.Set(user, db.COLLECTION_USER)
	if err != nil {
		revel.ERROR.Printf("Error while saving new user : %s", err.Error())
		return err
	}

	return nil
}

// RegisterEventDay registers a new event
func RegisterEventDay(event *model.EventDay) error {
	err := db.Database.Set(event, db.COLLECTION_EVENT_DAY)
	if err != nil {
		revel.ERROR.Printf("Error while saving new event : %s", err.Error())
		return err
	}

	return nil
}

// RegisterEventDay registers a new blacklist
func RegisterBlacklist(blacklist *model.Blacklist) error {
	err := db.Database.Set(blacklist, db.COLLECTION_BLACKLIST)
	if err != nil {
		revel.ERROR.Printf("Error while saving new blacklist : %s", err.Error())
		return err
	}

	return nil
}

// RankingEventNumber gets from the daily events, the ranking by number of events
func RankingEventNumber(params ...string) (MapReduceData, error) {
	var result MapReduceData

	var mapfunc string
	if len(params) == 1 {
		mapfunc = fmt.Sprintf("function() { if (this.type == '%s') { emit(this.user, 1) } }", params[0])
	} else {
		mapfunc = fmt.Sprintf("function() { if (this.type == '%s' && this.language == '%s') { emit(this.user, 1) } }", params[0], params[1])
	}

	_, err := db.Database.MapReduce(
		mapfunc,
		"function (key, values) { return Array.sum(values) }",
		db.COLLECTION_EVENT_DAY,
		&result,
	)

	fmt.Println(len(result))

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing event number : %s", err.Error())
		return nil, err
	}

	sort.Sort(result)

	// Checks for limit of events per day (to remove bots)
	if params[0] == "pushevent" {
		index := 0
		for key, value := range result {
			if value.Value < revel.Config.IntDefault("blacklist.limit", 200) {
				index = key
				break
			} else {
				RegisterBlacklist(model.NewBlacklist(value.Key, fmt.Sprintf("Number of events too big (%d)", value.Value)))
				db.Database.Remove(map[string]string{"user": value.Key}, db.COLLECTION_EVENT_DAY)
			}
		}

		result = result[index:len(result)]
	}

	return result, nil
}

// RankingEventExperience gets from the daily events, the ranking by experience
func RankingEventExperience(params ...string) (MapReduceData, error) {
	var result MapReduceData

	var mapfunc string
	if len(params) == 1 {
		mapfunc = fmt.Sprintf("function() { if (this.type == '%s') { emit(this.user, this.experience) } }", params[0])
	} else {
		mapfunc = fmt.Sprintf("function() { if (this.type == '%s' && this.language == '%s') { emit(this.user, this.experience) } }", params[0], params[1])
	}

	_, err := db.Database.MapReduce(
		mapfunc,
		"function (key, values) { return Array.sum(values) }",
		db.COLLECTION_EVENT_DAY,
		&result,
	)

	fmt.Println(len(result))

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing event experience : %s", err.Error())
		return nil, err
	}

	sort.Sort(result)
	return result, nil
}

// ClearEventDay removes all elements from events collection
func ClearEventDay() error {
	_, err := db.Database.ClearCollection(db.COLLECTION_EVENT_DAY)
	if err != nil {
		revel.ERROR.Printf("Error while clearing events collection : %s", err.Error())
		return err
	}

	return nil
}
