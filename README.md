# Validate

A validator for `http.Request`.

> Very rough, don't use in production, etc, etc.

## Use

1. Create a number of `Rule` items to check against. Each rule
has a parameter and a callback that determines if a Rule passes
or not.

```go
nameRequired := validate.Rule{
    Param: "name",
    Check: func(r *http.Request, param string) error {
        if _, exists := r.Form[param]; !exists {
            return fmt.Error("%s is required", param)
        }

        return nil
    }
}
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