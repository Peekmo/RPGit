package services

import (
	"RPGit/app/db"
	"RPGit/app/model"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

// Map reduce data (implements sort.Interface)
type MapReduceData []KeyValue

type KeyValue struct {
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

// IsFilled allows to know if the data has been filled in database or not
// It checks for the number of repositories
func IsFilled() bool {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	var repositories []*model.Repository

	data := db.Database.GetQuery(map[string]string{}, db.COLLECTION_REPOSITORY)
	data.Limit(1).All(&repositories)

	return (len(repositories) == 1)
}

// GetUser gets a user from the database
func GetUser(username string) *model.User {
	db.Database.InitSession()
	defer db.Database.Session.Close()

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
	db.Database.InitSession()
	defer db.Database.Session.Close()

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
	db.Database.InitSession()
	defer db.Database.Session.Close()

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
	db.Database.InitSession()
	defer db.Database.Session.Close()

	var repositories []*model.Repository

	data := db.Database.GetQuery(map[string]string{"owner": strings.ToLower(username)}, db.COLLECTION_REPOSITORY)
	data.All(&repositories)

	return repositories
}

// UpdateUser updates the given user from the database
func UpdateUser(user *model.User) error {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	err := db.Database.Update(user.Id, user, db.COLLECTION_USER)
	if err != nil {
		revel.ERROR.Printf("Error while updating data : %s", err.Error())
		return err
	}

	return nil
}

// UpdateRepository updates the repository in the database
func UpdateRepository(repository *model.Repository) error {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	err := db.Database.Update(repository.Id, repository, db.COLLECTION_REPOSITORY)
	if err != nil {
		revel.ERROR.Printf("Error while updating data : %s", err.Error())
		return err
	}

	return nil
}

// RegisterRepository register in database a new repository
func RegisterRepository(repository *model.Repository) error {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	err := db.Database.Set(repository, db.COLLECTION_REPOSITORY)
	if err != nil {
		revel.ERROR.Printf("Error while saving new repository : %s", err.Error())
		return err
	}

	return nil
}

// RegisterRepository register in database a new repository
func RegisterUser(user *model.User) error {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	err := db.Database.Set(user, db.COLLECTION_USER)
	if err != nil {
		revel.ERROR.Printf("Error while saving new user : %s", err.Error())
		return err
	}

	return nil
}

// RegisterEventDay registers a new event
func RegisterEventDay(event *model.EventDay) error {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	err := db.Database.Set(event, db.COLLECTION_EVENT_DAY)
	if err != nil {
		revel.ERROR.Printf("Error while saving new event : %s", err.Error())
		return err
	}

	return nil
}

// RegisterEventDay registers a new blacklist
func RegisterBlacklist(blacklist *model.Blacklist) error {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	err := db.Database.Set(blacklist, db.COLLECTION_BLACKLIST)
	if err != nil {
		revel.ERROR.Printf("Error while saving new blacklist : %s", err.Error())
		return err
	}

	return nil
}

// RankingExperience returns 50 first users sorted by experience
func RankingExperience() []KeyValue {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	var users []*model.User

	data := db.Database.GetQuery(map[string]string{}, db.COLLECTION_USER).Sort("-experience").Limit(50)
	data.All(&users)

	var formatted []KeyValue
	for _, user := range users {
		formatted = append(formatted, KeyValue{user.Username, user.Experience})
	}

	return formatted
}

// RankingExperienceLanguage returns 50 first users sorted by experience for the given language
func RankingExperienceLanguage(language string) (MapReduceData, error) {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	var result MapReduceData

	mapfunc := fmt.Sprintf("function() { for (var lang in this.languages) { if (this.languages[lang].name == '%s') { emit(this.username, this.languages[lang].experience); return; } } }", language)

	_, err := db.Database.MapReduce(
		mapfunc,
		"function (key, values) { return Array.sum(values) }",
		db.COLLECTION_USER,
		&result,
	)

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing experience language : %s", err.Error())
		return nil, err
	}

	sort.Sort(result)
	return result, nil
}

// RankingGlobalEventNumber gets from the daily events, the ranking by number of events
func RankingGlobalEventNumber(params ...string) (MapReduceData, error) {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	var result MapReduceData

	var mapfunc string
	if len(params) == 1 {
		mapfunc = fmt.Sprintf("function() { for (var lang in this.languages) { emit(this.username, this.languages[lang].events.pushes) } }")
	} else {
		mapfunc = fmt.Sprintf("function() { for (var lang in this.languages) { if (this.languages[lang].name == '%s') { emit(this.username, this.languages[lang].events.pushes); return; } } }", params[1])
	}

	_, err := db.Database.MapReduce(
		mapfunc,
		"function (key, values) { return Array.sum(values) }",
		db.COLLECTION_USER,
		&result,
	)

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing event number : %s", err.Error())
		return nil, err
	}

	sort.Sort(result)

	return result, nil
}

// RankingEventNumber gets from the daily events, the ranking by number of events
func RankingEventNumber(params ...string) (MapReduceData, error) {
	db.Database.InitSession()
	defer db.Database.Session.Close()

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
				db.Database.Remove(map[string]string{"_id": value.Key}, db.COLLECTION_USER)
			}
		}

		result = result[index:len(result)]
	}

	return result, nil
}

// RankingEventExperience gets from the daily events, the ranking by experience
func RankingEventExperience(params ...string) (MapReduceData, error) {
	db.Database.InitSession()
	defer db.Database.Session.Close()

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

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing event experience : %s", err.Error())
		return nil, err
	}

	sort.Sort(result)
	return result, nil
}

// RankingAllEventTotal returns the total daily events by language
func RankingAllEventTotal(typeEvent string) (MapReduceData, error) {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	var result MapReduceData

	var mapfunc string
	mapfunc = fmt.Sprintf("function() { if (this.type == '%s' && this.language != \"Unknown\") { emit(this.language, 1) } }", typeEvent)

	_, err := db.Database.MapReduce(
		mapfunc,
		"function (key, values) { return Array.sum(values) }",
		db.COLLECTION_EVENT_DAY,
		&result,
	)

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing event total : %s", err.Error())
		return nil, err
	}

	sort.Sort(result)
	return result, nil
}

// GetAllLanguages gets the list of languages with their number of pushes
func GetAllLanguages() (MapReduceData, error) {
	var result MapReduceData

	err := cache.Get("all_languages", &result)
	if err != nil {
		db.Database.InitSession()
		defer db.Database.Session.Close()

		mapfunc := "function() { this.languages.forEach(function(language) { if (language.name != \"Unknown\" && language.name != \"\") emit(language.name, language.events.pushes);});}"

		_, err := db.Database.MapReduce(
			mapfunc,
			"function (key, values) { return Array.sum(values) }",
			db.COLLECTION_USER,
			&result,
		)

		if err != nil {
			revel.ERROR.Printf("Error while mapreducing languages : %s", err.Error())
			return nil, err
		}

		sort.Sort(result)
		cache.Set("all_languages", result, cache.DEFAULT)
	}

	return result, nil
}

// ClearEventDay removes all elements from events collection
func ClearEventDay() error {
	db.Database.InitSession()
	defer db.Database.Session.Close()

	_, err := db.Database.ClearCollection(db.COLLECTION_EVENT_DAY)
	if err != nil {
		revel.ERROR.Printf("Error while clearing events collection : %s", err.Error())
		return err
	}

	return nil
}

// FetchAllRankingData builds the result by language or not
func FetchAllRankingData(typeEvent, language string, useCache bool) map[string](map[string]MapReduceData) {
	var data map[string](map[string]MapReduceData) = make(map[string](map[string]MapReduceData))

	err := cache.Get(fmt.Sprintf("ranking-home-%s-%s", typeEvent, strings.Join(strings.Split(language, " "), "")), &data)
	if err != nil || useCache == false {
		if err != nil {
			fmt.Print(err.Error())
		}

		var dailyNumber MapReduceData
		var dailyExperience MapReduceData
		var globalExperience []KeyValue
		var globalNumber MapReduceData

		if language != "" {
			dailyNumber, _ = RankingEventNumber(typeEvent, language)
			dailyExperience, _ = RankingEventExperience(typeEvent, language)
			globalExperience, _ = RankingExperienceLanguage(language)
			globalNumber, _ = RankingGlobalEventNumber(typeEvent, language)
		} else {
			dailyNumber, _ = RankingEventNumber(typeEvent)
			dailyExperience, _ = RankingEventExperience(typeEvent)
			globalExperience = RankingExperience()
			globalNumber, _ = RankingGlobalEventNumber(typeEvent)
		}

		dailyLanguage, _ := RankingAllEventTotal(typeEvent)
		globalLanguage, _ := GetAllLanguages()

		data["daily"] = map[string]MapReduceData{
			"number":     dailyNumber[0:int(math.Min(float64(len(dailyNumber)), float64(50)))],
			"experience": dailyExperience[0:int(math.Min(float64(len(dailyExperience)), float64(50)))],
			"language":   dailyLanguage[0:int(math.Min(float64(len(dailyLanguage)), float64(50)))],
		}

		data["global"] = map[string]MapReduceData{
			"number":     globalNumber[0:int(math.Min(float64(len(globalNumber)), float64(50)))],
			"experience": globalExperience[0:int(math.Min(float64(len(globalExperience)), float64(50)))],
			"language":   globalLanguage[0:int(math.Min(float64(len(dailyLanguage)), float64(50)))],
		}

		cache.Set(fmt.Sprintf("ranking-home-%s-%s", typeEvent, strings.Join(strings.Split(language, " "), "")), data, cache.DEFAULT)
	}

	return data
}

// ClearRankingCaches updates the in-memory cache
func ClearRankingCaches() {
	revel.INFO.Print("Clearing memory cache...")

	cache.Delete("all_languages")
	languages, _ := GetAllLanguages()
	FetchAllRankingData("pushevent", "", false)
	for _, language := range languages {
		FetchAllRankingData("pushevent", language.Key, false)
	}

	Ban("/")
}
