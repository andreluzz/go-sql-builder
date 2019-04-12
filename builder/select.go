package builder

// Select returns a statement with columns
func Select(columns ...string) *Statement {
	return &Statement{
		Type:    "select",
		Columns: columns,
	}
}

// From defines statement from table
func (s *Statement) From(table string) *Statement {
	s.Table = table
	return s
}
