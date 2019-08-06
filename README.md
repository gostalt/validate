# Validate

A validator for the Go `http.Request`.

Gostalt Validate provides an easy-to-use package for validating
form values in Go. With it, you construct a ruleset, and use it
to ensure that an incoming request is satisfied.

## Usage

Firstly, create a number of `Rule` items to check against. Each
rule has a parameter and a callback that is ran to determine if
the Rule passes or not.

Validate ships with a number of [built-in checks] that you can
use:

[built-in checks]: #Validators

```go
alphaForename := validate.Rule{
  Param: "forename",
  Check: validation.Alpha,
}

alphaSurname := validate.Rule{
  Param: "surname",
  Check: validation.Alpha,
}
```

However, you are free to create your own rules by passing a func
to the Check field. The func must be compatible with `CheckFunc`:

```go
type CheckFunc func(*http.Request, string, Options) error
```

Inside the function, you can use the http.Request to extract any
parameters you wish to check—the second argument to the function
is the parameter, so you can retrieve it dynamically from the
request:

```go
// Below is the built-in Integer validator.
var Integer CheckFunc = func(r *http.Request, param string, _ Options) error {
	fail, _ := regexp.MatchString(`[^0-9]+`, r.Form.Get(param))

	if fail {
		return fmt.Errorf("%s must be an integer", param)
	}

	return nil
}
```

Some validators utilise an `Options` map to provide dynamic checks:

```go
long := validate.Rule{
  Param: "name",
  Check: validate.MaxLength,
  // Max length is 5 characters
  Options: validate.Options{"length": 5},
}
```

You can use the Options map in your own custom validators using
the third argument to the CheckFunc:

```go
// Below is the built-in MaxLength validator.
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
```

When you have created all the validators you wish to use against
a request, call `validate.Check`. This accepts an http.Request
and any number of rules, and returns a `Message` and an error:

```go
msgs, err := validate.Check(r, alphaForename, alphaSurname)
```

It may be more succinct to create a variable or function that
contains the rules, and use the spread operator:

```go
msgs, err := validate.Check(r, rules()...)
```

When creating validation as part of a JSON API, you can easily
create a response using `validate.Response` (after checking
whether the validation as failed or not!). All you need to do
is pass in a ResponseWriter and the failed validation messages:

```go
validate.Response(w, msgs)
```

This will automatically write a `422 Unprocessable Entity` header
and a `Content-Type: application/json` header to the response,
then, the errors will be wrapped in an `error` object:

```json
{
  "errors": [
    {
      "forename": [
        "forename is required",
        "forename must be longer than 5 characters"
      ],
      "surname": [
        "surname must be longer than 5 characters"
      ],
      "dob": [
        "dob does not satisfy date format 2006-01-02"
      ]
    }
  ]
}
```

Of course, you can manually interact with the `msgs` variable
that is returned from the `Check` method if you need to carry
out additional logic or handling of failed validation.

## Validators

Validate offers a number of built-in validators:

### Alpha

Fails if the parameter contains any characters that aren't in the
alphabet.

### Alphanumeric

Fails if the parameter contains any characters that are not an
alphabet letter or a number.

### Boolean

Fails if the parameter is not a boolean. Because request params
do not have a type, this validator fails if the value of the
param is not equal to `"true"`, `"false"`, `"1"` or `"0"`.

### Date

Fails if the parameter is not a date. This validator uses the
built-in Go date formats to check if the passed value is a date.
You can also pass a slice of strings as a `formats` key in the
Options struct.

### DateFormat

Fails if the parameter does not matched the `format` key's value
passed in the Options struct.

### Email

Fails if the parameter is not an email address.

### Integer

Fails if the parameter is not an integer. Uses the regex `^[0-9]+`
to determine a pass.

### MaxLength

Fails if the length of the parameter (number of characters) is
greater than the `length` value passed to the Options struct.

### MinLength

Fails if the length of the parameter (number of characters) is
less than the `length` value passed to the Options struct.

### NotRegex

Fails if the parameter is satisfied by the regex `pattern` passed
to the Options struct.

### Regex

Fails if the parameter does not satisfy the regex `pattern` passed
to the Options struct.

### Required

Fails if the parameter does not exist in the request.

### RFC3339

Fails if the value passed does not match the date pattern defined
by RFC3339.
