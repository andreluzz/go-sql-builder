package builder

// Update returns a statement with columns
func Update(table string, columns ...string) *Statement {
	return &Statement{
		Type:    "update",
		Table:   table,
		Columns: columns,
	}
}
