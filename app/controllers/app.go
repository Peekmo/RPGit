package controllers

import (
	"RPGithub/api"
	"RPGithub/app/services"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

// GetUser returns the given user's data
func (c App) GetUser(username string) revel.Result {
	user := services.GetUser(username)
	if nil == user {
		return api.HttpException(c.Controller, 404, "User not found")
	}

	c.Response.Out.Header().Add("Cache-Control", "max-age=100000, public")
	return c.RenderJson(user)
}

// GetUserRepositories gets repositories from the given user
func (c App) GetUserRepositories(username string) revel.Result {
	user := services.GetUser(username)
	if nil == user {
		return api.HttpException(c.Controller, 404, "User not found")
	}

	repositories := services.GetUserRepositories(username)
	c.Response.Out.Header().Add("Cache-Control", "max-age=100000, public")
	return c.RenderJson(repositories)
}

// GetRankingTypeNumber returns number of events by type
func (c *App) GetRankingTypeNumber(typeEvent string) revel.Result {
	data, _ := services.RankingEventNumber(typeEvent)

	return c.RenderJson(data)
}

// GetRankingTypeExperience returns number of experiences points by user
func (c *App) GetRankingTypeExperience(typeEvent string) revel.Result {
	data, _ := services.RankingEventExperience(typeEvent)

	return c.RenderJson(data)
}

// GetRankingTypeNumberLanguage returns number of events by language, by user
func (c *App) GetRankingTypeNumberLanguage(typeEvent, language string) revel.Result {
	data, _ := services.RankingEventNumber(typeEvent, language)

	return c.RenderJson(data)
}

// GetRankingTypeExperienceLanguage returns the number of experiences points by user & language
func (c *App) GetRankingTypeExperienceLanguage(typeEvent, language string) revel.Result {
	data, _ := services.RankingEventExperience(typeEvent, language)

	return c.RenderJson(data)
}

// GetAllLanguages returns all languages with their number of events
func (c *App) GetAllLanguages() revel.Result {
	data, _ := services.GetAllLanguages()

	return c.RenderJson(data)
}
