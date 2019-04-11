package builder

// Delete returns a statement
func Delete(table string) *Statement {
	return &Statement{
		Type:  "delete",
		Table: table,
	}
}
