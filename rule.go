package validate

import (
	"fmt"
	"net/http"
	"regexp"
)

// Rule represents a check to run on a request.
type Rule struct {
	// Param is the field in the request to check.
	Param string
	// Check is a callback that is ran against the request.
	Check CheckFunc
}

type CheckFunc func(*http.Request, string) error

var (
	// Required returns an error if the parameter is not in the request.
	// Additional checks should be made to ensure it is not empty, etc.
	Required CheckFunc = func(r *http.Request, param string) error {
		if _, exists := r.Form[param]; !exists {
			return fmt.Errorf("%s is required", param)
		}

		return nil
	}

	// Alpha returns an error if the parameter contains any characters
	// that are not in the alphabet, represented by the regular
	// expression `[a-zA-Z]+`.
	Alpha CheckFunc = func(r *http.Request, param string) error {
		fail, _ := regexp.MatchString(`[^a-zA-Z]+`, r.Form.Get(param))

		if fail {
			return fmt.Errorf("%s must only contain alphabetical characters", param)
		}

		return nil
	}

	Alphanumeric CheckFunc = func(r *http.Request, param string) error {
		fail, _ := regexp.MatchString(`[^a-zA-Z0-9]+`, r.Form.Get(param))

		if fail {
			return fmt.Errorf("%s must only contain alphanumeric characters", param)
		}

		return nil
	}

	Boolean CheckFunc = func(r *http.Request, param string) error {
		value := r.Form.Get(param)

		if value == "true" || value == "false" || value == "1" || value == "0" {
			return nil
		}

		return fmt.Errorf("%s must be a boolean value", param)
	}
)
