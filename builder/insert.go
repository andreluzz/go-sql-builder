package builder

// Insert returns a statement
func Insert(table string, columns ...string) *Statement {
	return &Statement{
		Type:    "insert",
		Table:   table,
		Columns: columns,
	}
}

//Return include in insert statement the return columns
func (s *Statement) Return(columns ...string) *Statement {
	s.ReturnColumns = append(s.ReturnColumns, columns...)
	return s
}
