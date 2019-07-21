package validate

// Message represents a failed validation. It contains details
// of the param that failed, as well as the error message from
// the rule that caused it to fail.
type Message struct {
	Error string `json:"error"`
	Param string `json:"param"`
}
