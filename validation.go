package validate

import (
	"context"
	"encoding/json"
	"net/http"
)

// Validator is responsible for collecting an http.Request and
// a list of rules and checking that the rules are satisfied
// by the given request.
type Validator struct {
	request *http.Request
	Rules   []Rule
}

// Respond is a helper method that writes the errors to the given
// http.ResponseWriter. This also sets an appropriate HTTP header
// and sets the content-type to JSON.
func Respond(w http.ResponseWriter, m Message) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	eb := make(map[string]map[string][]string)
	eb["errors"] = m

	d, _ := json.Marshal(eb)
	w.Write(d)
}

// Check is an all-in-one method of creating and running a new
// Validator. If there is no logic around adding rules, it is
// the easiest way to run a Validator.
func Check(r *http.Request, rule ...Rule) (Message, error) {
	return Make(r, rule...).Run()
}

// Make creates a new Validator based on the request and rules
// passed into it. The rules argument is optional. Rules can
// be added by calling `Add` on the returned Validator.
func Make(r *http.Request, rule ...Rule) *Validator {
	if r.Form == nil {
		r.ParseForm()
	}

	return &Validator{
		request: r,
		Rules:   rule,
	}
}

// Run determines if the given rules are satisfied by the request.
func (v *Validator) Run() (Message, error) {
	if len(v.Rules) == 0 {
		return nil, EmptyRuleset
	}

	vm := make(Message)

	for _, rule := range v.Rules {
		if err := rule.Check(v.request, rule.Param, rule.Options); err != nil {
			vm[rule.Param] = append(vm[rule.Param], err.Error())
		}
	}

	if len(vm) > 0 {
		return vm, ValidationFailed
	}

	return nil, nil
}

// Add adds an additional set of rules to the Validator.
func (v *Validator) Add(rules ...Rule) {
	v.Rules = append(v.Rules, rules...)
}

type Bag string

const ErrorBag Bag = "errorbag"

func ErrorContext(r *http.Request, msgs Message) context.Context {
	return context.WithValue(r.Context(), ErrorBag, msgs)
}
