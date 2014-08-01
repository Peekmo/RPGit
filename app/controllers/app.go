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
