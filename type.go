package sqlw

import (
	"errors"
	"regexp"
	"strings"
)

// The following error is returned when the string validation of SQLQuery or SQLMutation fails.
var (
	ErrNotSQLQuery    = errors.New("it is not query statement")
	ErrNotSQLMutation = errors.New("it is not mutation statement")
)

var (
	mutationRe = regexp.MustCompile("[^insert.*$|^update.*$|^delete.*$]")
)

// SQLQuery provides for query(SELECT) statements extensions to string
type SQLQuery string

// Validate validates whether the string is a query(SELECT) statements.
func (s SQLQuery) Validate() error {
	str := strings.ToLower(string(s))
	if !strings.HasPrefix(str, "select") {
		return ErrNotSQLQuery
	}
	return nil
}

// String returns a transformed string.
func (s SQLQuery) String() string {
	return string(s)
}

// SQLMutation provides for mutation(INSERT|UPDATE|DELETE) statements extensions to string
type SQLMutation string

// Validate validates whether the string is a mutation(INSERT|UPDATE|DELETE) statements.
func (s SQLMutation) Validate() error {
	str := strings.ToLower(string(s))
	if !mutationRe.MatchString(str) {
		return ErrNotSQLMutation
	}
	return nil
}

// String returns a transformed string.
func (s SQLMutation) String() string {
	return string(s)
}
