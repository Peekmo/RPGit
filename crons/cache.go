package crons

import (
	"RPGithub/app/services"
)

// WarmCache warms in-memory cache and clear varnish
type WarmCache struct{}

// Run implements revel conrs.Job interface
func (this WarmCache) Run() {
	services.ClearRankingCaches()
}
