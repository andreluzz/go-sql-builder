package builder

//Select returns a statement with columns
func Select(columns ...string) *Statement {
	return &Statement{
		Type:    "select",
		Columns: columns,
	}
}

//From defines statement from table
func (s *Statement) From(table string) *Statement {
	s.Table = table
	return s
}

// Where adds a where condition.
// query can be Builder or string. value is used only if where type is string.
func (s *Statement) Where(where interface{}, value ...interface{}) *Statement {
	switch where := where.(type) {
	case string:
		s.WhereCond = append(s.WhereCond, Expr(where, value...))
	case Builder:
		s.WhereCond = append(s.WhereCond, where)
	}
	return s
}

// Join add foreng table with inner join
func (s *Statement) Join(table, on string) *Statement {
	s.JoinTable = append(s.JoinTable, join(table, on))
	return s
}
