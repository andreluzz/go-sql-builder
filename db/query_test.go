package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructSelectQuery(t *testing.T) {
	user := &User{
		ID: "059fa339-025c-4104-ab55-c764d3028f63",
	}
	query, values := StructSelectQuery("users", user)
	fmt.Println(values)
	assert.Equal(t, "query", query)
}
