package builder

import "strings"

// Statement represents a sql query
type Statement struct {
	Type      string
	Table     string
	Columns   []string
	WhereCond Builder
	JoinTable []Builder
	Data      []interface{}
}

// Values defines the input data to insert and update
func (s *Statement) Values(values ...interface{}) *Statement {
	s.Data = append(s.Data, values...)
	return s
}

// Prepare build the query that will be executed
func (s *Statement) Prepare(q Query) error {
	switch s.Type {
	case "select":
		return prepareSelect(s, q)
	case "insert":
		return prepareInsert(s, q)
	case "update":
		return prepareUpdate(s, q)
	case "delete":
		return prepareDelete(s, q)
	}

	return nil
}

func prepareSelect(s *Statement, q Query) error {
	q.WriteString("SELECT ")
	q.WriteString(strings.Join(s.Columns, ", "))
	q.WriteString(" FROM ")
	q.WriteString(s.Table)

	//joins
	if len(s.JoinTable) > 0 {
		for _, join := range s.JoinTable {
			err := join.Prepare(q)
			if err != nil {
				return err
			}
		}
	}

	err := s.WhereCond.Prepare(q)
	if err != nil {
		return err
	}

	return nil
}

func prepareInsert(s *Statement, q Query) error {
	q.WriteString("INSERT INTO ")
	q.WriteString(s.Table)
	q.WriteString(" (")
	q.WriteString(strings.Join(s.Columns, ", "))
	q.WriteString(") ")

	q.WriteString("VALUES ")
	records := len(s.Data) / len(s.Columns)
	for i := 0; i < records; i++ {
		if i > 0 {
			q.WriteString(", ")
		}
		q.WriteString("(")
		for i := 0; i < len(s.Columns); i++ {
			if i > 0 {
				q.WriteString(", ")
			}
			q.WriteString("?")
		}
		q.WriteString(")")
	}

	q.WriteValue(s.Data...)

	return nil
}

func prepareUpdate(s *Statement, q Query) error {
	q.WriteString("UPDATE ")
	q.WriteString(s.Table)
	q.WriteString(" SET ")

	for i, col := range s.Columns {
		if i > 0 {
			q.WriteString(", ")
		}
		q.WriteString(col)
		q.WriteString(" = ?")
	}

	q.WriteValue(s.Data...)

	err := s.WhereCond.Prepare(q)
	if err != nil {
		return err
	}

	return nil
}

func prepareDelete(s *Statement, q Query) error {
	q.WriteString("DELETE FROM ")
	q.WriteString(s.Table)

	err := s.WhereCond.Prepare(q)
	if err != nil {
		return err
	}

	return nil
}
