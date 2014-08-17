package controllers

import (
	"RPGithub/api"
	"RPGithub/app/services"
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
	return c.RenderJson(fetchAllRankingData(typeEvent, ""))
}

// GetHomeRankingsLanguage gets all ranking data for the given language)
func (c *App) GetHomeRankingsLanguage(typeEvent, language string) revel.Result {
	return c.RenderJson(fetchAllRankingData(typeEvent, language))
}

// GetAllLanguages returns all languages with their number of events
func (c *App) GetAllLanguages() revel.Result {
	data, _ := services.GetAllLanguages()

	return c.RenderJson(data)
}

// fetchAllRankingData builds the result by language or not
func fetchAllRankingData(typeEvent, language string) map[string](map[string]interface{}) {
	var data map[string](map[string]interface{}) = make(map[string](map[string]interface{}))
	var dailyNumber services.MapReduceData
	var dailyExperience services.MapReduceData

	if language != "" {
		dailyNumber, _ = services.RankingEventNumber(typeEvent, language)
		dailyExperience, _ = services.RankingEventExperience(typeEvent, language)
	} else {
		dailyNumber, _ = services.RankingEventNumber(typeEvent)
		dailyExperience, _ = services.RankingEventExperience(typeEvent)
	}

	dailyLanguage, _ := services.RankingAllEventTotal(typeEvent)

	data["daily"] = map[string]interface{}{
		"number":     dailyNumber[0:int(math.Min(float64(len(dailyNumber)), float64(50)))],
		"experience": dailyExperience[0:int(math.Min(float64(len(dailyExperience)), float64(50)))],
		"language":   dailyLanguage[0:int(math.Min(float64(len(dailyLanguage)), float64(50)))],
	}
	return data
}
