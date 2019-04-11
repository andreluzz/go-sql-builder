package builder

func join(table, on string) Builder {
	return PrepareFunc(func(q Query) error {
		q.WriteString(" JOIN ")
		q.WriteString(table)
		q.WriteString(" ON ")
		q.WriteString(on)
		return nil
	})
}
