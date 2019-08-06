package validate

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestCheckCreatesValidatorAndRunsIt(t *testing.T) {
	r, _ := http.NewRequest("GET", "localhost", nil)

	rule := Rule{
		Param: "forename",
		Check: func(r *http.Request, param string, _ Options) error {
			return errors.New("fail")
		},
	}

	rule2 := Rule{
		Param: "forename",
		Check: func(r *http.Request, param string, _ Options) error {
			return errors.New("fail")
		},
	}

	msgs, _ := Check(r, rule, rule2)
	if len(msgs) == 0 {
		fmt.Println("expected an error, didn't get one")
		t.FailNow()
	}
}

func TestMakeReturnsErrorWithNoRules(t *testing.T) {
	r, _ := http.NewRequest("GET", "localhost", nil)

	validator := Make(r)

	if _, err := validator.Run(); err == nil {
		fmt.Println("no error returned from empty validator")
		t.FailNow()
	}
}

func TestCanAddRules(t *testing.T) {
	r, _ := http.NewRequest("GET", "localhost", nil)

	validator := Make(r)

	rule := Rule{}

	validator.Add(rule)

	if len(validator.Rules) == 0 {
		fmt.Println("validator empty, expected 1 rule")
		t.FailNow()
	}
}

func TestFailureReturnsValidationMessage(t *testing.T) {
	r, _ := http.NewRequest("GET", "localhost", nil)

	validator := Make(r)

	rule := Rule{
		Param: "forename",
		Check: func(r *http.Request, param string, _ Options) error {
			return errors.New("fail")
		},
	}

	validator.Add(rule)

	messages, _ := validator.Run()

	if len(messages) == 0 {
		fmt.Println("expected a validation message, got none")
		t.FailNow()
	}
}

func TestValidatorRunsRuleCheck(t *testing.T) {
	r, _ := http.NewRequest("GET", "localhost", nil)

	validator := Make(r)

	rule := Rule{
		Param: "forename",
		Check: func(r *http.Request, param string, _ Options) error {
			return errors.New("forced failure")
		},
	}

	validator.Add(rule)

	if _, err := validator.Run(); err == nil {
		fmt.Println("expected validator to fail. It didn't.")
		t.FailNow()
	}
}

func TestValidatorReturnsRuleCheckErrorMessage(t *testing.T) {
	r, _ := http.NewRequest("GET", "localhost", nil)

	validator := Make(r)

	rule := Rule{
		Param: "forename",
		Check: func(r *http.Request, param string, _ Options) error {
			return errors.New("forced failure")
		},
	}

	validator.Add(rule)

	messages, _ := validator.Run()

	if messages["forename"][0] != "forced failure" {
		fmt.Println("expected the message to contain the rule error. It didn't.")
		t.FailNow()
	}
}

func TestValidatorReturnsParamInError(t *testing.T) {
	r, _ := http.NewRequest("GET", "localhost", nil)

	validator := Make(r)

	rule := Rule{
		Param: "forename",
		Check: func(r *http.Request, param string, _ Options) error {
			return errors.New("forced failure")
		},
	}

	validator.Add(rule)

	messages, _ := validator.Run()

	for param := range messages {
		if param != "forename" {
			fmt.Println("expected the message to contain the param. It didn't.")
			t.FailNow()
		}
	}
}

func TestParamHasNestedErrors(t *testing.T) {
	r, _ := http.NewRequest("GET", "localhost", nil)

	validator := Make(r)

	rule := Rule{
		Param: "forename",
		Check: func(r *http.Request, param string, _ Options) error {
			return errors.New("forced failure")
		},
	}

	rule2 := Rule{
		Param: "forename",
		Check: func(r *http.Request, param string, _ Options) error {
			return errors.New("forced failure")
		},
	}

	validator.Add(rule, rule2)

	messages, _ := validator.Run()

	for _, msgs := range messages {
		if len(msgs) != 2 {
			fmt.Println("expected the param to contain two errors. It didn't")
			t.FailNow()
		}
	}
}
