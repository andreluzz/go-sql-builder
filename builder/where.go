package builder

//Where add a new where clause to the query with all the conditions.
//Only uses values attribute if where is a string.
func where(where interface{}, values ...interface{}) Builder {
	return PrepareFunc(func(q Query) error {
		q.WriteString(" WHERE ")
		switch where := where.(type) {
		case string:
			q.WriteString("(")
			Raw(where, values...).Prepare(q)
			q.WriteString(")")
		case Builder:
			err := And(where).Prepare(q)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Where adds a where condition.
// query can be Builder or string. value is used only if where type is string.
func (s *Statement) Where(conditions interface{}, values ...interface{}) *Statement {
	s.WhereCond = where(conditions, values...)
	return s
}
