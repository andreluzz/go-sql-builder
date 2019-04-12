package builder

// Insert returns a statement
func Insert(table string, columns ...string) *Statement {
	return &Statement{
		Type:    "insert",
		Table:   table,
		Columns: columns,
	}
}
