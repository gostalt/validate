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
	// Options is a map that is passed to the check func.
	Options Options
}

// Options is a map of strings to values that can be used inside
// a CheckFunc to dynamically determine if a criteria has passed.
// Think max length, greater than, between, etc.
type Options map[string]interface{}

// CheckFunc is a function that uses the request, parameter and
// any passed options to determine if a Rule is passed.
type CheckFunc func(*http.Request, string, Options) error

var (
	// Required returns an error if the parameter is not in the request.
	// Additional checks should be made to ensure it is not empty, etc.
	Required CheckFunc = func(r *http.Request, param string, _ Options) error {
		if _, exists := r.Form[param]; !exists {
			return fmt.Errorf("%s is required", param)
		}

		return nil
	}

	// Alpha returns an error if the parameter contains any characters
	// that are not in the alphabet, represented by the regular
	// expression `[a-zA-Z]+`.
	Alpha CheckFunc = func(r *http.Request, param string, _ Options) error {
		fail, _ := regexp.MatchString(`[^a-zA-Z]+`, r.Form.Get(param))

		if fail {
			return fmt.Errorf("%s must only contain alphabetical characters", param)
		}

		return nil
	}

	// Alphanumeric returns an error if the parameter contains
	// any characters that are not letters or numbers.
	Alphanumeric CheckFunc = func(r *http.Request, param string, _ Options) error {
		fail, _ := regexp.MatchString(`[^a-zA-Z0-9]+`, r.Form.Get(param))

		if fail {
			return fmt.Errorf("%s must only contain alphanumeric characters", param)
		}

		return nil
	}

	// Integer returns an error if the parameter contains any
	// character that is not a digit.
	Integer CheckFunc = func(r *http.Request, param string, _ Options) error {
		fail, _ := regexp.MatchString(`[^0-9]+`, r.Form.Get(param))

		if fail {
			return fmt.Errorf("%s must be an integer", param)
		}

		return nil
	}

	// Boolean returns an error if the parameter contains a value
	// that is not boolean. Because these values are coming in
	// via a HTTP request (and are therefore strings), a boolean
	// value must be inferred.
	Boolean CheckFunc = func(r *http.Request, param string, _ Options) error {
		value := r.Form.Get(param)

		if value == "true" || value == "false" || value == "1" || value == "0" {
			return nil
		}

		return fmt.Errorf("%s must be a boolean value", param)
	}

	// MaxLength returns an error if the parameter length (number
	// of characters) exceeds the length set in the Options map
	// passed to the Rule.
	MaxLength CheckFunc = func(r *http.Request, param string, o Options) error {
		value := r.Form.Get(param)

		max, ok := o["length"].(int)
		if !ok {
			max = 0
		}

		if len(value) > max {
			return fmt.Errorf("%s cannot be longer than %d characters", param, max)
		}

		return nil
	}

	// MinLength returns an error if the parameter length (number
	// of characters) is shorter than the length set in the Options
	// map passed to the Rule.
	MinLength CheckFunc = func(r *http.Request, param string, o Options) error {
		value := r.Form.Get(param)

		min, ok := o["length"].(int)
		if !ok {
			min = 0
		}

		if len(value) < min {
			return fmt.Errorf("%s must be longer than %d characters", param, min)
		}

		return nil
	}

	// Regex returns an error if the parameter does not satisfy
	// the regular expression passed in the Options map.
	Regex CheckFunc = func(r *http.Request, param string, o Options) error {
		value := r.Form.Get(param)

		pattern, ok := o["pattern"].(string)
		if !ok {
			return fmt.Errorf("unable to create regex to validate %s parameter", param)
		}

		if pass, _ := regexp.MatchString(pattern, value); !pass {
			return fmt.Errorf("%s did not match regex `%s`", param, pattern)
		}

		return nil
	}
)
