package model

import (
	"strings"
)

type User struct {
	Id            string          `json:"id" bson:"_id"`
	Name          string          `json:"name"`
	Username      string          `json:"username"`
	Gravatar      int             `json:"gravatar"`
	Avatar        string          `json:"avatar"`
	Languages     []*Language     `json:"languages"`
	Organizations []*Organization `json:"organizations"`
}

type Language struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	RecvEvents *Events `json:"events_received"` // Events received on a personal repo
	SentEvents *Events `json:"events_sent"`     // Events sent on a personal or other repo
	Experience int     `json:"experience"`
}

type Events struct {
	Pushes       int `json:"pushes"`
	Pullrequests int `json:"pullrequests"`
	Issues       int `json:"issues"`
	Forks        int `json:"forks"`
	Watches      int `json:"watches"`
	Stars        int `json:"stars"`
}

type Organization struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
	Events  *Events  `json:"events"`
}

type Repository struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Url          string `json:"url"`
	Language     string `json:"language"`
	Owner        string `json:"owner"`
	Organization bool   `json:"organization"`
}

// NewUser creates a new user
func NewUser(username string) *User {
	return &User{Id: strings.ToLower(username), Username: username}
}

// NewUser creates a new user
func NewRepository(id int, name string, size int, url, language, owner string, organization bool) *Repository {
	return &Repository{id, name, size, url, language, owner, organization}
}

// GetLanguage gets an instance of the given language id
// It creates it if it does not exists
func (this *User) GetLanguage(id string) *Language {
	for _, language := range this.Languages {
		if language.Id == id {
			return language
		}
	}

	language := &Language{id, id, &Events{}, &Events{}, 0}
	this.Languages = append(this.Languages, language)

	return language
}
