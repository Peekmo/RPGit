package services

import (
	"net/http"

	"github.com/revel/revel"
)

// Ban ban the given path
func Ban(path string) {
	if revel.Config.BoolDefault("varnish.enabled", false) {
		ip := revel.Config.StringDefault("varnish.ip", "127.0.0.1")
		revel.INFO.Printf("%s%s", ip, path)

		client := &http.Client{}
		req, err := http.NewRequest("BAN", "http://"+ip+path, nil)
		if err != nil {
			revel.ERROR.Printf("Error on BAN request creation on %s : %s", path, err.Error())
		}

		_, err = client.Do(req)
		if err != nil {
			revel.ERROR.Printf("Error on BAN request on %s : %s", path, err.Error())
		}
	}
}
