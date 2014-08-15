package controllers

import (
	"RPGithub/app/services"

	"github.com/revel/revel"
)

// Interface is the controller for different views
type Interface struct {
	*revel.Controller
}

// Index renders the home page
func (c Interface) Index() revel.Result {
	return c.Render()
}

// User renders user profile
func (c Interface) User(username string) revel.Result {
	user := services.GetUser(username)
	if nil == user {
		return c.NotFound("User not found")
	}

	return c.Render(user)
}
