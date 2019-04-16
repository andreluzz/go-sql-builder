package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructSelectQuery(t *testing.T) {
	users := []User{
		{
			FirstName: "teste 1",
			LastName:  "teste 1",
			Email:     "ahahhaha",
		},
		{
			FirstName: "teste 2",
			LastName:  "teste 2",
		},
	}

	query, values := StructMultipleInsertQuery("users", users)
	fmt.Println(values)
	assert.Equal(t, "query", query)
}
