package services

import (
	"html/template"

	"github.com/revel/revel"
)

// RegisterHelpers is a public which registers all template helpers
func RegisterHelpers() {
	revel.TemplateFuncs["colorLanguage"] = ColorLanguages
	revel.TemplateFuncs["invertColorLanguage"] = InvertColorLanguages
	revel.TemplateFuncs["importDate"] = GetImportDate
}

// ColorLanguages prints the color of the event number label in user profile
func ColorLanguages(current, teal, green int) template.HTML {
	if current >= green {
		return template.HTML("green")
	} else if current >= teal {
		return template.HTML("teal")
	} else {
		return template.HTML("red")
	}
}

// InvertColorLanguages prints the color of the event number label in user profile
func InvertColorLanguages(current, teal, red int) template.HTML {
	if current >= red {
		return template.HTML("red")
	} else if current >= teal {
		return template.HTML("teal")
	} else {
		return template.HTML("green")
	}
}

// GetImportDate returns the date of the first import setted in config
func GetImportDate() template.HTML {
	return template.HTML(revel.Config.StringDefault("imports.begin", "01/01/1970"))
}
