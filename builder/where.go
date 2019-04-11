package builder

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
