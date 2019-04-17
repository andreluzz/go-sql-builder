package db

import (
	"fmt"
	"testing"

	"github.com/andreluzz/go-sql-builder/builder"
	"github.com/stretchr/testify/assert"
)

func TestStructSelectQuery(t *testing.T) {
	users := []User{}

	query, values, err := StructSelectQuery("users", &users, builder.Equal("id", "0000001"))
	assert.NoError(t, err, "invalid interface")
	fmt.Println(values)
	assert.Equal(t, "query", query)
}
