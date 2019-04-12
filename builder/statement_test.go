package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectStatement(t *testing.T) {
	statement := Select("tn.column", "tn.column2", "tx.column").From("table_name tn").Join("table_external tx", "tx.fk_id = tn.id").Where(
		Or(
			And(
				NotEqual("tn.column", nil),
				Equal("tn.column2", "teste"),
				GreaterOrEqual("tn.column3", 10),
				LowerOrEqual("tn.column4", 10),
				GreaterThen("tn.column5", 10),
				LowerThen("tn.column6", 10),
			),
			And(
				NotEqual("tn.column", "200"),
				Equal("tn.column2", nil),
			),
		),
	)

	query := NewQuery()
	statement.Prepare(query)

	assert.Equal(t, "SELECT tn.column, tn.column2, tx.column FROM table_name tn JOIN table_external tx ON tx.fk_id = tn.id WHERE (((tn.column IS NOT NULL) AND (tn.column2 = ?) AND (tn.column3 >= ?) AND (tn.column4 <= ?) AND (tn.column5 > ?) AND (tn.column6 < ?)) OR ((tn.column != ?) AND (tn.column2 IS NULL)))", query.String())
	var list []interface{}
	list = append(list, "teste", 10, 10, 10, 10, "200")
	assert.Equal(t, list, query.Value())
}

func TestInsertStatement(t *testing.T) {
	statement := Insert("table_name", "column", "column2").Values("Apartamento", 1000).Values("Casa", 750)

	query := NewQuery()
	statement.Prepare(query)

	assert.Equal(t, "INSERT INTO table_name (column, column2) VALUES (?, ?), (?, ?)", query.String())
	var list []interface{}
	list = append(list, "Apartamento", 1000, "Casa", 750)
	assert.Equal(t, list, query.Value())
}

func TestUpdateStatement(t *testing.T) {
	var list []interface{}
	list = append(list, "Apartamento", 1000)

	statement := Update("table_name", "column", "column2").Values(list...).Where("id = ?", 999)

	list = append(list, 999)

	query := NewQuery()
	statement.Prepare(query)

	assert.Equal(t, "UPDATE table_name SET column = ?, column2 = ? WHERE (id = ?)", query.String())
	assert.Equal(t, list, query.Value())
}

func TestDeleteStatement(t *testing.T) {
	statement := Delete("table_name tn").Where(
		Equal("tn.column", 100),
	)

	query := NewQuery()
	statement.Prepare(query)

	assert.Equal(t, "DELETE FROM table_name tn WHERE (tn.column = ?)", query.String())
	var list []interface{}
	list = append(list, 100)
	assert.Equal(t, list, query.Value())
}
