/**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /*

  @author Axel Anceau - 2014
  Package api contains general tools

*/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/ /**/
package api

import (
	"github.com/revel/revel"
)

var (
	HttpMessages map[int]string
)

// ApiError data structure which contains the error
type ApiError struct {
	Error *ApiDataError `json:"error"` // Error's data
}

// ApiDataError is the structure which defines the error message
type ApiDataError struct {
	Code    int    `json:"code"`    // Status code
	Text    string `json:"text"`    // Http string value
	Message string `json:"message"` // Custom message
}

// Response returns a revel response in API's format
func Response(c *revel.Controller, object interface{}, statusCode int) revel.Result {
	c.Response.WriteHeader(statusCode, "application/json")

	return c.RenderJson(object)
}

// HttpException returns a Response with the given status and message.
// The corresponding HttpMessage is added
func HttpException(c *revel.Controller, statusCode int, message string) revel.Result {
	return Response(c, ApiError{&ApiDataError{statusCode, HttpMessages[statusCode], message}}, statusCode)
}

// init function - Builds HttpMessages's map
func init() {
	HttpMessages = make(map[int]string)
	HttpMessages = map[int]string{
		100: "Continue",
		101: "Switching Protocols",
		102: "Processing",
		118: "Connection timed out",

		200: "OK",
		201: "Created",
		202: "Accepted",
		203: "Non-Authoritative Information",
		204: "No Content",
		205: "Reset Content",
		206: "Partial Content",
		207: "Multi-Status",
		210: "Content Different",
		226: "IM Used",

		300: "Multiple Choices",
		301: "Moved Permanently",
		302: "Moved Temporarily",
		303: "See Other",
		304: "Not Modified",
		305: "Use Proxy",
		306: "(none)",
		307: "Temporary Redirect",
		308: "Permanent Redirect",
		310: "Too many Redirects",

		400: "Bad Request",
		401: "Unauthorized",
		402: "Payment Required",
		403: "Forbidden",
		404: "Not Found",
		405: "Method Not Allowed",
		406: "Not Acceptable",
		407: "Proxy Authentication Required",
		408: "Request Time-out",
		409: "Conflict",
		410: "Gone",
		411: "Length Required",
		412: "Precondition Failed",
		413: "Request Entity Too Large",
		414: "Request-URI Too Long",
		415: "Unsupported Media Type",
		416: "Request range unsatisfiable",
		417: "Expectation failed",
		418: "I'm a teapot",
		422: "Unprocessable entity",
		423: "Locked",
		424: "Method failure",
		425: "Unordered Collection",
		426: "Upgrade Required",
		428: "Precondition Required",
		429: "Too Many Requests",
		431: "Request Header Fields Too Large",
		449: "Retry With",
		450: "Blocked by Windows Parent Controls",
		456: "Unrecoverable Error",
		499: "Client has closed connection",

		500: "Internal Server Error",
		501: "Not Implemented",
		502: "Bad Gateway",
		503: "Service Unavailable",
		504: "Gateway Time-out",
		505: "HTTP Version not supported",
		506: "Variant also negociate",
		507: "Insufficient storage",
		508: "Loop detected",
		509: "Bandwidth Limit Exceeded",
		510: "Not extended",
		520: "Web serer is returning an unknown error",
	}
}
