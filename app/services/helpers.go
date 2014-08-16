package services

import (
	"html/template"

	"github.com/revel/revel"
)

// RegisterHelpers is a public which registers all template helpers
func RegisterHelpers() {
	revel.TemplateFuncs["colorLanguage"] = ColorLanguages
}

// colorLanguages prints the color of the event number label in user profile
func ColorLanguages(current, teal, green int) template.HTML {
	if current >= green {
		return template.HTML("green")
	} else if current >= teal {
		return template.HTML("teal")
	} else {
		return template.HTML("red")
	}
}
