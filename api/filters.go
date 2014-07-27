/**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /*

  @author Axel Anceau - 2014
  Package api contains general tools

*/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/
package api

import (
	"fmt"
	"github.com/revel/revel"
	"runtime/debug"
)

/**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /*

  PanicFilter renders a panic as JSON

  @see revel/panic.go

*/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/
func PanicFilter(c *revel.Controller, fc []revel.Filter) {
	defer func() {
		if err := recover(); err != nil && err != "HttpException" {
			error := revel.NewErrorFromPanic(err)
			if error == nil {
				revel.ERROR.Print(err, "\n", string(debug.Stack()))
				c.Response.Out.WriteHeader(500)
				c.Response.Out.Write(debug.Stack())
				return
			}

			revel.ERROR.Print(err, "\n", error.Stack)
			c.Result = HttpException(c, 500, fmt.Sprint(err))
		}
	}()
	fc[0](c, fc[1:])
}
