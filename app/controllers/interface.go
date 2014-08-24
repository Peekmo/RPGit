package controllers

import (
	"RPGit/app/model"
	"RPGit/app/services"
	"sort"

	"github.com/revel/revel"
)

// Interface is the controller for different views
type Interface struct {
	*revel.Controller
}

// Language is a language structure for the menu
type Language struct {
	Key   string
	Value int
}

// Index renders the home page
func (c Interface) Index() revel.Result {
	languages, _ := services.GetAllLanguages()

	var total int
	for _, value := range languages {
		total += value.Value
	}

	for key, value := range languages {
		languages[key].Value = (value.Value * 1000000) / total
	}

	return c.Render(languages)
}

// User renders user profile
func (c Interface) User(username string) revel.Result {
	user := services.GetUser(username)
	if nil == user {
		return c.NotFound("User not found")
	}

	sort.Sort(sort.Reverse(model.LanguageArray(user.Languages)))
	sort.Sort(sort.Reverse(model.RepositoryArray(user.Repositories)))

	var languages []*Language
	for _, language := range user.Languages {
		languages = append(languages, &Language{language.Name, language.Level})
	}

	var showValueLanguage = true
	return c.Render(user, languages, showValueLanguage)
}
