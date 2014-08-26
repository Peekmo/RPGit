package services

import (
	"RPGit/app/db"
	"RPGit/app/model"
	"fmt"
	"math"
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
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	var users []*model.User

	data := database.GetQuery(map[string]string{}, db.COLLECTION_USER)
	data.Limit(1).All(&users)

	return (len(users) == 1)
}

// GetUser gets a user from the database
func GetUser(username string) *model.User {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	var user *model.User

	userData := database.Get(strings.ToLower(username), db.COLLECTION_USER)
	err := userData.One(&user)
	if err != nil {
		return nil
	}

	return user
}

// IsBlacklisted checks if the given user is blacklisted or not
func IsBlacklisted(name string) bool {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	var blacklist *model.Blacklist

	blacklistData := database.Get(strings.ToLower(name), db.COLLECTION_BLACKLIST)
	err := blacklistData.One(&blacklist)
	if err != nil {
		return false
	}

	return true
}

// UpdateUser updates the given user from the database
func UpdateUser(user *model.User) error {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	err := database.Update(user.Id, user, db.COLLECTION_USER)
	if err != nil {
		revel.ERROR.Printf("Error while updating data : %s", err.Error())
		return err
	}

	return nil
}

// RegisterRepository register in database a new repository
func RegisterUser(user *model.User) error {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	err := database.Set(user, db.COLLECTION_USER)
	if err != nil {
		revel.ERROR.Printf("Error while saving new user : %s", err.Error())
		return err
	}

	database.Index(db.COLLECTION_USER, "username")
	database.Index(db.COLLECTION_USER, "languages.name")

	return nil
}

// RegisterEventDay registers a new event
func RegisterEventDay(event *model.EventDay) error {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	err := database.Set(event, db.COLLECTION_EVENT_DAY)
	if err != nil {
		revel.ERROR.Printf("Error while saving new event : %s", err.Error())
		return err
	}

	database.Index(db.COLLECTION_EVENT_DAY, "user")
	database.Index(db.COLLECTION_EVENT_DAY, "language")
	database.Index(db.COLLECTION_EVENT_DAY, "type")

	return nil
}

// RegisterEventDay registers a new blacklist
func RegisterBlacklist(blacklist *model.Blacklist) error {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	err := database.Set(blacklist, db.COLLECTION_BLACKLIST)
	if err != nil {
		revel.ERROR.Printf("Error while saving new blacklist : %s", err.Error())
		return err
	}

	return nil
}

// RankingExperience returns 50 first users sorted by experience
func RankingExperience() []KeyValue {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	var users []*model.User

	data := database.GetQuery(map[string]string{}, db.COLLECTION_USER).Sort("-experience").Limit(50)
	data.All(&users)

	var formatted []KeyValue
	for _, user := range users {
		formatted = append(formatted, KeyValue{user.Username, user.Experience})
	}

	return formatted
}

// RankingExperienceLanguage returns 50 first users sorted by experience for the given language
func RankingExperienceLanguage(language string) (MapReduceData, error) {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	var result MapReduceData

	mapfunc := fmt.Sprintf("function() { for (var lang in this.languages) { if (this.languages[lang].name == '%s') { emit(this.username, this.languages[lang].experience); return; } } }", language)

	_, err := database.MapReduce(
		mapfunc,
		"function (key, values) { return Array.sum(values) }",
		"",
		db.COLLECTION_USER,
		map[string]string{"languages.name": language},
	)

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing experience language : %s", err.Error())
		return nil, err
	}

	data := database.GetQuery(map[string]string{}, db.COLLECTION_MAPREDUCE).Sort("-value")
	data.All(&result)

	return result, nil
}

// RankingGlobalEventNumber gets from the daily events, the ranking by number of events
func RankingGlobalEventNumber(params ...string) (MapReduceData, error) {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	var result MapReduceData

	var mapfunc string
	if len(params) == 1 {
		mapfunc = fmt.Sprintf("function() { for (var lang in this.languages) { emit(this.username, this.languages[lang].events.pushes) } }")
	} else {
		mapfunc = fmt.Sprintf("function() { for (var lang in this.languages) { if (this.languages[lang].name == '%s') { emit(this.username, this.languages[lang].events.pushes); return; } } }", params[1])
	}

	_, err := database.MapReduce(
		mapfunc,
		"function (key, values) { return Array.sum(values) }",
		"",
		db.COLLECTION_USER,
		map[string]string{},
	)

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing event number : %s", err.Error())
		return nil, err
	}

	data := database.GetQuery(map[string]string{}, db.COLLECTION_MAPREDUCE).Sort("-value").Limit(50)
	data.All(&result)

	return result, nil
}

// RankingEventNumber gets from the daily events, the ranking by number of events
func RankingEventNumber(params ...string) (MapReduceData, error) {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	var result MapReduceData
	var query = make(map[string]string)

	if len(params) == 1 {
		query["type"] = params[0]
	} else {
		query["type"] = params[0]
		query["language"] = params[1]
	}

	_, err := database.MapReduce(
		"function() { emit(this.user, 1) }",
		"function (key, values) { return Array.sum(values) }",
		"user",
		db.COLLECTION_EVENT_DAY,
		query,
	)

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing event number : %s", err.Error())
		return nil, err
	}

	data := database.GetQuery(map[string]string{}, db.COLLECTION_MAPREDUCE).Sort("-value")
	data.All(&result)

	// Checks for limit of events per day (to remove bots)
	if params[0] == "pushevent" {
		index := 0
		for key, value := range result {
			if value.Value < revel.Config.IntDefault("blacklist.limit", 250) {
				index = key
				break
			} else {
				RegisterBlacklist(model.NewBlacklist(value.Key, fmt.Sprintf("Number of events too big (%d)", value.Value)))

				database.Remove(map[string]string{"user": value.Key}, db.COLLECTION_EVENT_DAY)
				database.Remove(map[string]string{"_id": value.Key}, db.COLLECTION_USER)
			}
		}

		result = result[index:len(result)]
	}

	return result, nil
}

// RankingEventExperience gets from the daily events, the ranking by experience
func RankingEventExperience(params ...string) (MapReduceData, error) {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	var result MapReduceData
	var query = make(map[string]string)

	if len(params) == 1 {
		query["language"] = params[0]
	} else if len(params) == 2 {
		query["language"] = params[0]
		query["type"] = params[1]
	}

	_, err := database.MapReduce(
		"function() { emit(this.user, this.experience) }",
		"function (key, values) { return Array.sum(values) }",
		"user",
		db.COLLECTION_EVENT_DAY,
		query,
	)

	if err != nil {
		revel.ERROR.Printf("Error while mapreducing event experience : %s", err.Error())
		return nil, err
	}

	data := database.GetQuery(map[string]string{}, db.COLLECTION_MAPREDUCE).Sort("-value").Limit(50)
	data.All(&result)

	return result, nil
}

// RankingAllEventTotal returns the total daily events by language
func RankingAllEventTotal(typeEvent string) (MapReduceData, error) {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	var result MapReduceData

	err := cache.Get("all_languages_daily", &result)
	if err != nil {
		_, err := database.MapReduce(
			"function() { emit(this.language, 1);}",
			"function (key, values) { return Array.sum(values) }",
			"language",
			db.COLLECTION_EVENT_DAY,
			map[string](interface{}){"type": typeEvent, "language": map[string]string{"$ne": "Unknown"}},
		)

		if err != nil {
			revel.ERROR.Printf("Error while mapreducing event total : %s", err.Error())
			return nil, err
		}

		data := database.GetQuery(map[string]string{}, db.COLLECTION_MAPREDUCE).Sort("-value").Limit(50)
		data.All(&result)

		cache.Set("all_languages_daily", result, cache.DEFAULT)
	}

	return result, nil
}

// GetAllLanguages gets the list of languages with their number of pushes
func GetAllLanguages() (MapReduceData, error) {
	var result MapReduceData

	err := cache.Get("all_languages", &result)
	if err != nil {
		session := db.Database.InitSession()
		database := db.Database.Copy(session)
		defer session.Close()

		mapfunc := "function() { this.languages.forEach(function(language) { if (language.name != \"Unknown\") { emit(language.name, language.events.pushes);}});}"

		_, err := database.MapReduce(
			mapfunc,
			"function (key, values) { return Array.sum(values) }",
			"",
			db.COLLECTION_USER,
			map[string]string{},
		)

		if err != nil {
			revel.ERROR.Printf("Error while mapreducing languages : %s", err.Error())
			return nil, err
		}

		data := database.GetQuery(map[string]string{}, db.COLLECTION_MAPREDUCE).Sort("-value")
		data.All(&result)

		cache.Set("all_languages", result, cache.DEFAULT)
	}

	return result, nil
}

// ClearEventDay removes all elements from events collection
func ClearEventDay() error {
	session := db.Database.InitSession()
	database := db.Database.Copy(session)
	defer session.Close()

	_, err := database.ClearCollection(db.COLLECTION_EVENT_DAY)
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
			dailyExperience, _ = RankingEventExperience(language)
			globalExperience, _ = RankingExperienceLanguage(language)
			globalNumber, _ = RankingGlobalEventNumber(typeEvent, language)
		} else {
			dailyNumber, _ = RankingEventNumber(typeEvent)
			dailyExperience, _ = RankingEventExperience()
			globalExperience = RankingExperience()
			globalNumber, _ = RankingGlobalEventNumber(typeEvent)
		}

		dailyLanguage, _ := RankingAllEventTotal(typeEvent)
		globalLanguage, _ := GetAllLanguages()

		data["daily"] = map[string]MapReduceData{
			"number":     dailyNumber[0:int(math.Min(float64(len(dailyNumber)), float64(50)))],
			"experience": dailyExperience,
			"language":   dailyLanguage,
		}

		data["global"] = map[string]MapReduceData{
			"number":     globalNumber,
			"experience": globalExperience,
			"language":   globalLanguage,
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

	cache.Delete("all_languages_daily")
	RankingAllEventTotal("pushevent")

	FetchAllRankingData("pushevent", "", false)
	var total = len(languages)
	for key, language := range languages {
		revel.WARN.Printf("Language %s (%d/%d)", language.Key, key+1, total)
		FetchAllRankingData("pushevent", language.Key, false)
	}

	Ban("/")
}
