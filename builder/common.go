package builder

// Values defines the input data to insert and update
func (s *Statement) Values(values ...interface{}) *Statement {
	s.Data = append(s.Data, values...)
	return s
}

func join(table, on string) Builder {
	return PrepareFunc(func(q Query) error {
		q.WriteString(" JOIN ")
		q.WriteString(table)
		q.WriteString(" ON ")
		q.WriteString(on)
		return nil
	})
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
