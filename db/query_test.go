package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructUpdateQuery(t *testing.T) {
	user := &User{
		ID:    "00000001",
		Email: "user@teste.com",
	}
	query, value, _ := StructUpdateQuery("users", user, "email")
	fmt.Println(value)
	assert.Equal(t, "query", query)
}
