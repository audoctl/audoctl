// Package docs Audoctl API
//
// Documentation for Audoctl
//
//	BasePath: /
//	Version: 1.0.0
//	Title: Audoctl Swagger
//	Consumes:
//	- application/json
//	Produces:
//	- application/json
//
// swagger:meta
package docs

//go:generate swag init --parseDependency --parseInternal

// swagger:response ErrorResponse
type ErrorResponseWrapper struct {
	// in: body
	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}
