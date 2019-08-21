package validate

type Error string

const (
	EmptyRuleset     Error = "attempted to run a validator with an empty rule set"
	ValidationFailed Error = "validation failed"
)

func (e Error) Error() string {
	return string(e)
}
