package builder

func buildCond(q Query, pred string, cond ...Builder) error {
	for i, c := range cond {
		if i > 0 {
			q.WriteString(" ")
			q.WriteString(pred)
			q.WriteString(" ")
		}
		q.WriteString("(")
		err := c.Prepare(q)
		if err != nil {
			return err
		}
		q.WriteString(")")
	}
	return nil
}

// And creates AND from a list of conditions.
func And(cond ...Builder) Builder {
	return PrepareFunc(func(query Query) error {
		return buildCond(query, "AND", cond...)
	})
}

// Or creates OR from a list of conditions.
func Or(cond ...Builder) Builder {
	return PrepareFunc(func(query Query) error {
		return buildCond(query, "OR", cond...)
	})
}

//GreaterThen greater then condition
func GreaterThen(column string, value interface{}) Builder {
	return PrepareFunc(func(q Query) error {
		q.WriteString(column)
		q.WriteString(" > ?")
		q.WriteValue(value)
		return nil
	})
}

//LowerThen lower then condition
func LowerThen(column string, value interface{}) Builder {
	return PrepareFunc(func(q Query) error {
		q.WriteString(column)
		q.WriteString(" < ?")
		q.WriteValue(value)
		return nil
	})
}

//GreaterOrEqual greater or equal condition
func GreaterOrEqual(column string, value interface{}) Builder {
	return PrepareFunc(func(q Query) error {
		q.WriteString(column)
		q.WriteString(" >= ?")
		q.WriteValue(value)
		return nil
	})
}

//LowerOrEqual lower or equal condition
func LowerOrEqual(column string, value interface{}) Builder {
	return PrepareFunc(func(q Query) error {
		q.WriteString(column)
		q.WriteString(" <= ?")
		q.WriteValue(value)
		return nil
	})
}

// Equal crates a equal comparison
func Equal(column string, value interface{}) Builder {
	return PrepareFunc(func(q Query) error {
		if value == nil {
			q.WriteString(column)
			q.WriteString(" IS NULL")
			return nil
		}
		q.WriteString(column)
		q.WriteString(" = ?")
		q.WriteValue(value)
		return nil
	})
}

//NotEqual creates a not equal comparison
func NotEqual(column string, value interface{}) Builder {
	return PrepareFunc(func(q Query) error {
		if value == nil {
			q.WriteString(column)
			q.WriteString(" IS NOT NULL")
			return nil
		}
		q.WriteString(column)
		q.WriteString(" != ?")
		q.WriteValue(value)
		return nil
	})
}
