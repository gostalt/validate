package validate

import (
	"net/http"
)

// Rule represents a check to run on a request.
type Rule struct {
	// Param is the field in the Request to check.
	Param string
	// Check is a callback that is ran on the request. The
	// second argument is the param to check. If the check
	// fails, an error should be returned with details why.
	Check func(*http.Request, string) error
}
