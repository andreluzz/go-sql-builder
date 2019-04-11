package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectStatement(t *testing.T) {
	statement := Select("tn.column", "tn.column2", "tx.column").From("table_name tn").Join("table_external tx", "tx.fk_id = tn.id").Where(
		And(
			Eq("tn.column", "100"),
			Eq("tn.column2", "teste"),
		),
	)

	query := NewQuery()
	statement.Prepare(query)

	assert.Equal(t, "SELECT tn.column, tn.column2, tx.column FROM table_name tn JOIN table_external tx ON tx.fk_id = tn.id WHERE ((tn.column = ?) AND (tn.column2 = ?))", query.String())
}

func TestInsertStatement(t *testing.T) {
	statement := Insert("table_name", "column", "column2").Values("Apartamento", 1000).Values("Casa", 750)

	query := NewQuery()
	statement.Prepare(query)

	assert.Equal(t, "INSERT INTO table_name (column, column2) VALUES (?, ?), (?, ?)", query.String())
}

func TestUpdateStatement(t *testing.T) {
	statement := Update("table_name", "column", "column2").Values("Apartamento", 1000).Where("id = ?", 1000)

	query := NewQuery()
	statement.Prepare(query)

	assert.Equal(t, "UPDATE table_name SET column = ?, column2 = ? WHERE (id = ?)", query.String())
}
