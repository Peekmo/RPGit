package controllers

import (
	// "RPGithub/app/services"
	"github.com/revel/revel"
)

type Interface struct {
	*revel.Controller
}

// Index
func (c Interface) Index() revel.Result {
	return c.Render()
}
