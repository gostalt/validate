# Validate

A validator for `http.Request`.

> Very rough, don't use in production, etc, etc.

## Use

1. Create a number of `Rule` items to check against. Each rule
has a parameter and a callback that determines if a Rule passes
or not.

There are a number of built-in checks that you can use, rather than
reinventing the wheel:

```go
nameIsAlpha := validate.Rule{
  Param: "name",
  Check: validate.Alpha,
}
```

However, you are free to create your own rules by passing a function
to the Check field. This function must take an http.Request and a string
and return an error:

```go
nameRequired := validate.Rule{
    Param: "name",
    Check: func(r *http.Request, param string) error {
        if _, exists := r.Form[param]; !exists {
            return fmt.Error("%s is required", param)
        }

        return nil
    },
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

When using built-in validation rules, you can inline the struct
to make it more readable at a glance:

```go
validate.Rule{"name", validate.MaxLength, validate.Options{"length": 5}}
```

2. Call `validate.Check`. This accepts an `http.Request` and any
number of Rules, and returns a slice of `Message`s and an error:

```go
msgs, err := validate.Check(r, nameRequired)
```

3. To automatically handle writing a response, you can call the
`validate.Response` method, and pass it an `http.ResponseWriter`
and the slice of `Message`s:

```go
validate.Response(w, msgs)
```

This will automatically write a `422` header and a
`Content-Type: application/json` header to the response, then,
the errors will be wrapped in an `error` object:

```json
{
  "errors": [
    {
      "error": "name is required",
      "param": "name"
    }
  ]
}
```

Of course, you can manually interact with the `msgs` variable
that is returned from the `Check` method if you need to carry
out additional logic or handling of failed validation.

## Available Validators

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
