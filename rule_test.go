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
			[]string{"abc", "1.5", ""},
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
		{
			NotRegex,
			[]string{"letters99", "__1", "66666__"},
			[]string{"55555aa", "514tomy", "1810Lucy"},
			Options{"pattern": `[0-9]+[a-zA-Z]+`},
		},
		{
			Email,
			[]string{"me@tomm.us", "me+99__.asd@subdomain.tomm.us"},
			[]string{"me@something@tomm.us", "juststring", "me space@tomm.us"},
			nil,
		},
		{
			MXEmail,
			[]string{"me@tomm.us", "lucyduggleby@hotmail.co.uk"},
			[]string{"me@something@addasadsdn2343567hgbf.com", "juststring", "me@space@tomm.us"},
			nil,
		},
		{
			RFC3339,
			[]string{"1993-10-18T10:10:10Z", "1992-06-22T10:10:10-05:00", "2006-01-02T15:04:05+01:00"},
			[]string{"1993-10-18", "1992-06-22"},
			nil,
		},
		{
			RFC1123,
			[]string{"Tue, 22 Jun 1992 10:00:00 GMT", "Tue, 18 Oct 1993 10:00:00 GMT"},
			[]string{"1993-10-18", "1992-06-22"},
			nil,
		},
		{
			RFC822,
			[]string{"22 Jun 92 10:00 GMT", "18 Oct 93 13:00 GMT"},
			[]string{"1992-06-22"},
			nil,
		},
		{
			UnixDate,
			[]string{"Mon Jan 16 15:04:05 MST 2006", "Tue Jun 22 10:00:00 GMT 1992"},
			[]string{"1993-10-18", "1990-11-11"},
			nil,
		},
		{
			DateFormat,
			[]string{"2016/02/29", "2019/10/18", "1992/06/22"},
			[]string{"2016-02-29", "2019-10-18", "1992-06-22"},
			Options{"format": "2006/01/02"},
		},
		{
			Date,
			[]string{"1993-10-18T10:10:10-02:00", "22 Jun 92 15:04 UTC", "2019-08-01"},
			[]string{"2016/02/29", "Monday 02 Jan 2006"},
			Options{"formats": []string{"2006-01-02"}},
		},
		{
			Date,
			[]string{"1993-10-18T10:10:10-02:00", "22 Jun 92 15:04 UTC", "2019-08-01"},
			[]string{"2016/02/29", "Monday 02 Jan 2006"},
			Options{"formats": []string{"2006-01-02"}},
		},
	}

	for _, rule := range rules {
		// First, ensure the check passes
		for _, value := range rule.Passes {
			r.Form.Set("parameter", value)
			msgs, _ := Check(r, Rule{"parameter", rule.Check, rule.Options})
			if len(msgs) > 0 {
				fmt.Println("Got an error, expected none:", msgs["parameter"])
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

func BenchmarkInteger(b *testing.B) {
	r, _ := http.NewRequest("GET", "localhost", nil)
	for n := 0; n < b.N; n++ {
		Integer(r, "example", nil)
	}
}

func BenchmarkDate(b *testing.B) {
	r, _ := http.NewRequest("GET", "localhost", nil)
	for n := 0; n < b.N; n++ {
		Date(r, "example", nil)
	}
}

func BenchmarkRFC3339(b *testing.B) {
	r, _ := http.NewRequest("GET", "localhost", nil)
	for n := 0; n < b.N; n++ {
		RFC3339(r, "example", nil)
	}
}

func BenchmarkEmail(b *testing.B) {
	r, _ := http.NewRequest("GET", "localhost", nil)
	for n := 0; n < b.N; n++ {
		Email(r, "example", nil)
	}
}
