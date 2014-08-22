package controllers

import (
	"RPGit/api"
	"RPGit/app/services"
	"math"

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

	max := math.Min(float64(len(data)), float64(50))

	return c.RenderJson(data[0:int(max)])
}

// GetRankingTypeExperience returns number of experiences points by user
func (c *App) GetRankingTypeExperience(typeEvent string) revel.Result {
	data, _ := services.RankingEventExperience(typeEvent)

	max := math.Min(float64(len(data)), float64(50))

	return c.RenderJson(data[0:int(max)])
}

// GetRankingTypeNumberLanguage returns number of events by language, by user
func (c *App) GetRankingTypeNumberLanguage(typeEvent, language string) revel.Result {
	data, _ := services.RankingEventNumber(typeEvent, language)

	max := math.Min(float64(len(data)), float64(50))

	return c.RenderJson(data[0:int(max)])
}

// GetRankingTypeExperienceLanguage returns the number of experiences points by user & language
func (c *App) GetRankingTypeExperienceLanguage(typeEvent, language string) revel.Result {
	data, _ := services.RankingEventExperience(typeEvent, language)

	max := math.Min(float64(len(data)), float64(50))

	return c.RenderJson(data[0:int(max)])
}

// GetRankingAllTypeTotal returns the total daily events by language
func (c *App) GetRankingAllTypeTotal(typeEvent string) revel.Result {
	data, _ := services.RankingAllEventTotal(typeEvent)

	max := math.Min(float64(len(data)), float64(50))

	return c.RenderJson(data[0:int(max)])
}

// GetHomeRankings gets all ranking data (no matter the language)
func (c *App) GetHomeRankings(typeEvent string) revel.Result {
	return c.RenderJson(services.FetchAllRankingData(typeEvent, "", true))
}

// GetHomeRankingsLanguage gets all ranking data for the given language)
func (c *App) GetHomeRankingsLanguage(typeEvent, language string) revel.Result {
	return c.RenderJson(services.FetchAllRankingData(typeEvent, language, true))
}

// GetAllLanguages returns all languages with their number of events
func (c *App) GetAllLanguages() revel.Result {
	data, _ := services.GetAllLanguages()

	return c.RenderJson(data)
}
