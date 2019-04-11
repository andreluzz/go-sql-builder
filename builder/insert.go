package builder

//Insert returns a statement
func Insert(table string, columns ...string) *Statement {
	return &Statement{
		Type:    "insert",
		Table:   table,
		Columns: columns,
	}
}

//Values defines the input data to insert and update
func (s *Statement) Values(values ...interface{}) *Statement {
	s.Data = append(s.Data, values...)
	return s
}
