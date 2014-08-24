package app

import (
	"RPGit/api"
	"RPGit/app/services"
	"RPGit/crons"

	"github.com/revel/revel"
	"github.com/revel/revel/modules/jobs/app/jobs"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		api.PanicFilter,               // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register custom template helpers
	services.RegisterHelpers()

	// register startup functions with OnAppStart
	revel.OnAppStart(func() {
		// Init database
		services.InitDatabase()

		// Defines CRONS
		jobs.Schedule("cron.import", crons.Import{})
		jobs.Now(crons.WarmCache{})
		jobs.Now(crons.FullImport{})
	})
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Cache-Control", "s-maxage=100000, private")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
