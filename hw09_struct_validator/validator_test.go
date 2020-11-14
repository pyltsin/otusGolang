package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in         interface{}
		countError int
	}{
		{App{Version: "123"}, 1},
		{App{Version: "12345"}, 0},
		{Response{Code: 100, Body: ""}, 1},
		{Response{Code: 200, Body: ""}, 0},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 20, Email: "test@test.ru", Role: "admin", Phones: []string{"12345678901"}, meta: nil}, 0},
		{User{ID: "12345678901234567890123456789012345", Name: "", Age: 20, Email: "test@test.ru", Role: "admin", Phones: []string{"12345678901"}, meta: nil}, 1},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 5, Email: "test@test.ru", Role: "admin", Phones: []string{"12345678901"}, meta: nil}, 1},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 20, Email: "test", Role: "admin", Phones: []string{"12345678901"}, meta: nil}, 1},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 20, Email: "test@test.ru", Role: "", Phones: []string{"12345678901"}, meta: nil}, 1},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 20, Email: "test@test.ru", Role: "admin", Phones: []string{"1234567890"}, meta: nil}, 1},
		{User{ID: "12345678901234567890123456789012345", Name: "", Age: 5, Email: "test", Role: "1", Phones: []string{"1234567890"}, meta: nil}, 5},
		{User{ID: "12345678901234567890123456789012345", Name: "", Age: 20, Email: "test", Role: "1", Phones: []string{"1234567890"}, meta: nil}, 4},
		{User{ID: "12345678901234567890123456789012345", Name: "", Age: 90, Email: "test", Role: "1", Phones: []string{"1234567890"}, meta: nil}, 5},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			validate := Validate(tt.in)
			assert.Equal(t, tt.countError, len(validate))
		})
	}
}
