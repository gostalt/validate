package validate

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// Required returns an error if the parameter is not in the request.
// Additional checks should be made to ensure it is not empty, etc.
var Required CheckFunc = func(r *http.Request, param string, _ Options) error {
	if _, exists := r.Form[param]; !exists {
		return fmt.Errorf("%s is required", param)
	}

	return nil
}

// Alpha returns an error if the parameter contains any characters
// that are not in the alphabet, represented by the regular
// expression `[a-zA-Z]+`.
var Alpha CheckFunc = func(r *http.Request, param string, _ Options) error {
	fail, _ := regexp.MatchString(`[^a-zA-Z]+`, r.Form.Get(param))

	if fail {
		return fmt.Errorf("%s must only contain alphabetical characters", param)
	}

	return nil
}

// Alphanumeric returns an error if the parameter contains
// any characters that are not letters or numbers.
var Alphanumeric CheckFunc = func(r *http.Request, param string, _ Options) error {
	fail, _ := regexp.MatchString(`[^a-zA-Z0-9]+`, r.Form.Get(param))

	if fail {
		return fmt.Errorf("%s must only contain alphanumeric characters", param)
	}

	return nil
}

// Integer returns an error if the parameter cannot be converted
// to an integer.
var Integer CheckFunc = func(r *http.Request, param string, _ Options) error {
	_, err := strconv.Atoi(r.Form.Get(param))
	if err != nil {
		return fmt.Errorf("%s must be an integer", param)
	}

	return nil
}

// Boolean returns an error if the parameter contains a value
// that is not boolean. Because these values are coming in
// via a HTTP request (and are therefore strings), a boolean
// value must be inferred.
var Boolean CheckFunc = func(r *http.Request, param string, _ Options) error {
	value := r.Form.Get(param)

	if value == "true" || value == "false" || value == "1" || value == "0" {
		return nil
	}

	return fmt.Errorf("%s must be a boolean value", param)
}

// MaxLength returns an error if the parameter length (number
// of characters) exceeds the length set in the Options map
// passed to the Rule.
var MaxLength CheckFunc = func(r *http.Request, param string, o Options) error {
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
var MinLength CheckFunc = func(r *http.Request, param string, o Options) error {
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
var Regex CheckFunc = func(r *http.Request, param string, o Options) error {
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

// NotRegex returns an error if the parameter value is satisfied
// by the regular expression passed in the Options map.
var NotRegex CheckFunc = func(r *http.Request, param string, o Options) error {
	value := r.Form.Get(param)

	pattern, ok := o["pattern"].(string)
	if !ok {
		return fmt.Errorf("unable to create regex to validate %s parameter", param)
	}

	if pass, _ := regexp.MatchString(pattern, value); pass {
		return fmt.Errorf("%s must not match regex `%s`", param, pattern)
	}

	return nil
}

// MXEmail looks up the MX Records on a domain to check if a record exists. If
// an MX record exists, it is likely that the email address is real. This is
// smarter than just checking if an email address fits a certain format.
var MXEmail CheckFunc = func(r *http.Request, param string, o Options) error {
	if err := Email(r, param, nil); err != nil {
		return err
	}

	timeout, ok := o["timeout"].(int)
	if !ok {
		timeout = 5
	}

	domain := getDomain(r.Form.Get(param))
	records, err := getMXRecords(r.Context(), domain, timeout)
	if err != nil {
		return fmt.Errorf("the host %s is not a valid email provider", domain)
	}

	if len(records) == 0 {
		return fmt.Errorf("no MX records exist for %s", param)
	}

	return nil
}

// Email returns an error if the parameter value is not a valid
// email address.
var Email CheckFunc = func(r *http.Request, param string, _ Options) error {
	value := r.Form.Get(param)

	atCount := strings.Count(value, "@")

	// If there is not one @ sign in the string, it is not
	// a valid email address.
	if atCount != 1 {
		return fmt.Errorf("%s is not a valid email address", param)
	}

	if pass, _ := regexp.MatchString(`^[^@\s]+@[^@\s]+$`, value); pass {
		return nil
	}

	return fmt.Errorf("%s is not a valid email address", param)
}

// RFC3339 returns an error if the parameter does not satisfy
// the RFC3339 format.
var RFC3339 CheckFunc = func(r *http.Request, param string, _ Options) error {
	return DateFormat(r, param, Options{"format": time.RFC3339})
}

// RFC1123 returns an error if the parameter does not satisfy
// the RFC1123 format.
var RFC1123 CheckFunc = func(r *http.Request, param string, _ Options) error {
	return DateFormat(r, param, Options{"format": time.RFC1123})
}

// RFC822 returns an error if the parameter does not satisfy the
// RFC822 format.
var RFC822 CheckFunc = func(r *http.Request, param string, _ Options) error {
	return DateFormat(r, param, Options{"format": time.RFC822})
}

// UnixDate returns an error if the parameter does not satisfy
// the format defined in Go's UnixDate const.
var UnixDate CheckFunc = func(r *http.Request, param string, _ Options) error {
	return DateFormat(r, param, Options{"format": time.UnixDate})
}

// DateFormat returns an error if the parameter does not
// satisfy the date format passed in the Options struct.
var DateFormat CheckFunc = func(r *http.Request, param string, o Options) error {
	value := r.Form.Get(param)

	format, ok := o["format"].(string)
	if !ok {
		return fmt.Errorf("unable to create date format string")
	}

	if _, err := time.Parse(format, value); err != nil {
		return fmt.Errorf("%s does not satisfy date format %s", param, format)
	}

	return nil
}

// Date is a comprehensive validator that returns an error if
// the parameter does not satisfy any of Go's built-in date
// formats.
//
// To validate against additional custom formats, you can pass
// a slice of strings to the Options struct using a `formats` key.
var Date CheckFunc = func(r *http.Request, param string, o Options) error {
	formats := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,

		// TODO: These are times... Maybe move them?
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	customFormats, exists := o["formats"]
	if exists {
		customFormats, ok := customFormats.([]string)
		if !ok {
			return fmt.Errorf("unable to create date format string")
		}

		formats = append(formats, customFormats...)
	}

	for _, format := range formats {
		if err := DateFormat(r, param, Options{"format": format}); err == nil {
			return nil
		}
	}

	return fmt.Errorf("%s does not satisfy and date format", param)
}

func getMXRecords(ctx context.Context, domain string, timeout int) ([]*net.MX, error) {
	rsv := net.Resolver{}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	return rsv.LookupMX(ctx, domain)
}

func getDomain(email string) string {
	parts := strings.Split(email, "@")
	return parts[len(parts)-1]
}
