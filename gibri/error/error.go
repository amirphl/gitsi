package error

import (
	"fmt"
)

type Scope int

const (
	Session Scope = iota
	System
)

type Error struct {
	scope  Scope
	detail string
}

func (j *Error) ShouldRetry() bool {
	return true
}

func (j *Error) String() string {
	return fmt.Sprintf("Error: %v %s\n", j.scope, j.detail)
}

func (j *Error) Error() string {
	return j.String()
}

func New(scope Scope, detail string) *Error {
	return &Error{
		scope:  scope,
		detail: detail,
	}
}
