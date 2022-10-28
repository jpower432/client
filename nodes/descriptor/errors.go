package descriptor

import "fmt"

type ParseError struct {
	Key string
	Err error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parse property key %q: %v", e.Key, e.Err)
}
