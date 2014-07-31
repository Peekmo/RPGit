package controllers

import (
	"RPGithub/api"
	"RPGithub/app/db"
	"RPGithub/app/model"
	"github.com/revel/revel"
	"strings"
)

type App struct {
	*revel.Controller
}

// GetUser returns the given user's data
func (c App) GetUser(username string) revel.Result {
	var user *model.User

	userData := db.Database.Get(strings.ToLower(username), db.COLLECTION_USER)
	err := userData.One(&user)
	if err != nil {
		return api.HttpException(c.Controller, 404, "User not found")
	}

	return c.RenderJson(user)
}
