package models

type User struct {
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	Username      string         `json:"username"`
	Gravatar      int            `json:"gravatar"`
	Avatar        string         `json:"avatar"`
	Languages     []Language     `json:"languages"`
	Organizations []Organization `json:"organizations"`
	Repositories  []Repository   `json:"repositories"`
}

type Language struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	RecvEvents []Events `json:"events_received"` // Events received on a personal repo
	SentEvents []Events `json:"events_sent"`     // Events sent on a personal or other repo
}

type Events struct {
	Pushes       int `json:"count"`
	Pullrequests int `json:"pullrequests"`
	Issues       int `json:"issues"`
	Forks        int `json:"forks"`
	Watches      int `json:"watches"`
	Stars        int `json:"stars"`
}

type Organization struct {
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Members []int    `json:"members"`
	Events  []Events `json:"events"`
}

type Repository struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	Url      string `json:"url"`
	Language string `json:"language"`
}
