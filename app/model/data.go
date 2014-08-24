package model

import (
	"strings"
)

type User struct {
	Id           string        `json:"id" bson:"_id"`
	Name         string        `json:"name"`
	Username     string        `json:"username"`
	Avatar       string        `json:"avatar"`
	Level        int           `json:"level"`
	Experience   int           `json:"experience"`
	Percent      int           `json:"percent"`
	Languages    []*Language   `json:"languages"`
	Repositories []*Repository `json:"repositories"`
}

type Language struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Events     *Events `json:"events"`
	Level      int     `json:"level"`
	Experience int     `json:"experience"`
	Percent    int     `json:"percent"`
}

type Events struct {
	Pushes       int `json:"pushes"`
	Pullrequests int `json:"pullrequests"`
	Issues       int `json:"issues"`
	Comments     int `json:"comments"`
	Forks        int `json:"forks"`
	Watches      int `json:"watches"`
	Stars        int `json:"stars"`
	Creates      int `json:"creates"`
	Deletes      int `json:"deletes"`
}

type Organization struct {
	Id      int      `json:"id" bson:"_id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
	Events  *Events  `json:"events"`
}

type Repository struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Url          string `json:"url"`
	Description  string `json:"description"`
	Language     string `json:"language"`
	Owner        string `json:"owner"`
	Organization string `json:"organization"`
	Wiki         bool   `json:"wiki"`
	Downloads    bool   `json:"downloads"`
	Forks        int    `json:"forks"`
	Stars        int    `json:"stars"`
	Issues       int    `json:"issues"`
	IsFork       bool   `json:"is_fork"`
}

type EventDay struct {
	Type       string
	Language   string
	User       string
	Experience int
	Repository int
	Date       string
}

type Blacklist struct {
	Name   string `json:"name" bson:"_id"`
	Reason string `json:"reason"`
}

// NewUser creates a new user
func NewUser(username string) *User {
	return &User{Id: strings.ToLower(username), Username: username}
}

// NewUser creates a new user
func NewRepository(id int, name string) *Repository {
	return &Repository{Id: id, Name: name}
}

// NewEventDay creates a new event day object
func NewEventDay(typeEvent, language, user string, experience, repository int, date string) *EventDay {
	return &EventDay{typeEvent, language, strings.ToLower(user), experience, repository, date}
}

// NewBlacklist creates a new blacklist
func NewBlacklist(name, reason string) *Blacklist {
	return &Blacklist{name, reason}
}

// AddExperience adds experience to the given user
// It also calculates its current level
func (this *User) AddExperience(xp int) {
	this.Experience += xp
	this.Level = getLevel(this.Level, this.Experience)
	this.Percent = getPercent(this.Level, this.Experience)
}

// GetLanguage gets an instance of the given language id
// It creates it if it does not exists
func (this *User) GetLanguage(id string) *Language {
	for _, language := range this.Languages {
		if language.Id == id {
			return language
		}
	}

	language := &Language{id, id, &Events{}, 1, 0, 0}
	this.Languages = append(this.Languages, language)

	return language
}

// GetRepository returns an instance of the repository
// It creates it if it does not already exists
func (this *User) GetRepository(id int, name string) *Repository {
	for _, repository := range this.Repositories {
		if repository.Id == id {
			return repository
		}
	}

	repository := NewRepository(id, name)
	this.Repositories = append(this.Repositories, repository)

	return repository
}

// AddExperience adds the given experience to the language
// It also calculates its current level
func (this *Language) AddExperience(xp int) {
	this.Experience += xp
	this.Level = getLevel(this.Level, this.Experience)
	this.Percent = getPercent(this.Level, this.Experience)
}

// ---- Tools

// getLevel returns the level for the given experience
func getLevel(level, experience int) int {
	for experience >= getExperience(level) {
		level += 1
	}

	return level
}

// getPercent returns the percentage of the current level
func getPercent(level, experience int) int {
	return experience * 100 / getExperience(level)
}

// getExperience returns experience needed for the next level
func getExperience(level int) int {
	return level * level * level * 2
}

// ----- Sorts
type LanguageArray []*Language
type RepositoryArray []*Repository

// Len returns array's length
func (this LanguageArray) Len() int {
	return len(this)
}

// Swap changes array's elements
func (this LanguageArray) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// Less returns if the given language as more or less experience than the other
func (this LanguageArray) Less(i, j int) bool {
	return this[i].Experience < this[j].Experience
}

// Len returns array's length
func (this RepositoryArray) Len() int {
	return len(this)
}

// Swap changes array's elements
func (this RepositoryArray) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// Less returns if the given repository as more or less stars than the other
func (this RepositoryArray) Less(i, j int) bool {
	return this[i].Stars < this[j].Stars
}
