package builder

type raw struct {
	Query string
	Value []interface{}
}

// Raw allows a manually created query to be used when current SQL syntax is not supported
func Raw(query string, value ...interface{}) Builder {
	return &raw{Query: query, Value: value}
}

// Prepare build the query that will be executed
func (raw *raw) Prepare(q Query) error {
	q.WriteString(raw.Query)
	q.WriteValue(raw.Value...)
	return nil
}
