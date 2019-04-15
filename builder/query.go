package builder

import (
	"strings"
)

// Builder defines a interface to
type Builder interface {
	Prepare(Query) error
}

// PrepareFunc function to assemble queries
type PrepareFunc func(Query) error

// Prepare creates the chain to build the query
func (p PrepareFunc) Prepare(query Query) error {
	return p(query)
}

// Query with string builder
type Query interface {
	WriteString(string) (int, error)
	String() string
	Reset()

	WriteValue(v ...interface{}) (err error)
	Value() []interface{}
}

type query struct {
	strings.Builder
	v []interface{}
}

// NewQuery creates a new Query.
func NewQuery() Query {
	return &query{}
}

// WriteValue populates the value array
func (q *query) WriteValue(v ...interface{}) error {
	q.v = append(q.v, v...)
	return nil
}

// Value returns the value array
func (q *query) Value() []interface{} {
	return q.v
}
