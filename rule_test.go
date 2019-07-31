package validate

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRules(t *testing.T) {
	// Make a dummy request variable.
	r, _ := http.NewRequest("POST", "localhost", nil)
	r.ParseForm()

	rules := []struct {
		Check  CheckFunc
		Passes []string
		Fails  []string
	}{
		{
			Alpha,
			[]string{"Alphabet", "lowercase", "UPPERCASE"},
			[]string{"Alphab3tic4l", "13567", "letters-and-dashes"},
		},
		{
			Alphanumeric,
			[]string{"Alphanumeric123", "123alpha", "123", "abc"},
			[]string{"number-letter-dash", "__", "--"},
		},
		{
			Boolean,
			[]string{"true", "false", "1", "0"},
			[]string{"2", "truthy", "falsy"},
		},
		{
			Integer,
			[]string{"123", "1", "0", "99"},
			[]string{"abc", "1.5"},
		},
	}

	for _, rule := range rules {
		// First, ensure the check passes
		for _, value := range rule.Passes {
			r.Form.Set("parameter", value)
			msgs, _ := Check(r, Rule{"parameter", rule.Check})
			if len(msgs) > 0 {
				fmt.Println("Got an error, expected none:", msgs[0].Error)
				fmt.Println("Value was", value)
				t.FailNow()
			}
		}

		// Then, ensure that it can fail
		for _, value := range rule.Fails {
			r.Form.Set("parameter", value)
			msgs, _ := Check(r, Rule{"parameter", rule.Check})
			if len(msgs) == 0 {
				fmt.Println("Expected an error, didn't get one")
				fmt.Println("Value was", value)
				t.FailNow()
			}
		}
	}
}

func TestRequiredRule(t *testing.T) {
	r, _ := http.NewRequest("GET", "localhost", nil)

	rule := Rule{
		Param: "anything",
		Check: Required,
	}

	msgs, _ := Check(r, rule)
	if len(msgs) == 0 {
		fmt.Println("expected an error, didn't get one")
		t.FailNow()
	}
}
