package db

import (
	"fmt"
	"testing"

	"github.com/andreluzz/go-sql-builder/builder"
	"github.com/stretchr/testify/assert"
)

func TestStructSelectQuery(t *testing.T) {
	group := Group{}

	query, values, err := StructSelectQuery("core_groups", &group, builder.Equal("core_groups.id", "0000001"))
	assert.NoError(t, err, "invalid interface")
	fmt.Println(values)
	assert.Equal(t, "query", query)
}
