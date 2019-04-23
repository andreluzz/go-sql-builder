package db

import (
	"fmt"
	"testing"

	"github.com/andreluzz/go-sql-builder/builder"
	"github.com/stretchr/testify/assert"
)

func TestStructSelectQuery(t *testing.T) {
	user := User{}

	query, values, err := StructSelectQuery("core_users", &user, builder.Equal("core_users.id", "57a97aaf-16da-44ef-a8be-b1caf52becd6"))
	assert.NoError(t, err, "invalid interface")
	fmt.Println(values)
	assert.Equal(t, "query", query)
}
