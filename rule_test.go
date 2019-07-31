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
		Check   CheckFunc
		Passes  []string
		Fails   []string
		Options Options
	}{
		{
			Alpha,
			[]string{"Alphabet", "lowercase", "UPPERCASE"},
			[]string{"Alphab3tic4l", "13567", "letters-and-dashes"},
			nil,
		},
		{
			Alphanumeric,
			[]string{"Alphanumeric123", "123alpha", "123", "abc"},
			[]string{"number-letter-dash", "__", "--"},
			nil,
		},
		{
			Boolean,
			[]string{"true", "false", "1", "0"},
			[]string{"2", "truthy", "falsy"},
			nil,
		},
		{
			Integer,
			[]string{"123", "1", "0", "99"},
			[]string{"abc", "1.5"},
			nil,
		},
		{
			MaxLength,
			[]string{"aaaa", "1111", "true", "----"},
			[]string{"too long by half", "TWO WEEEEEKS", "1111111"},
			Options{"length": 5},
		},
		{
			MinLength,
			[]string{"ok", "ye", "zz"},
			[]string{"a", "1", "_", "-"},
			Options{"length": 2},
		},
		{
			Regex,
			[]string{"55555aa", "514tomy", "1810Lucy"},
			[]string{"letters99", "__1", "66666__"},
			Options{"pattern": `[0-9]+[a-zA-Z]+`},
		},
	}

	for _, rule := range rules {
		// First, ensure the check passes
		for _, value := range rule.Passes {
			r.Form.Set("parameter", value)
			msgs, _ := Check(r, Rule{"parameter", rule.Check, rule.Options})
			if len(msgs) > 0 {
				fmt.Println("Got an error, expected none:", msgs[0].Error)
				fmt.Println("Value was", value)
				t.FailNow()
			}
		}

		// Then, ensure that it can fail
		for _, value := range rule.Fails {
			r.Form.Set("parameter", value)
			msgs, _ := Check(r, Rule{"parameter", rule.Check, rule.Options})
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
