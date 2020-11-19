package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"errors"
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

	DefectIn struct {
		Code int `validate:"in:a,b,c"`
	}
	DefectLen struct {
		Code string `validate:"len:a"`
	}
	EmptyValidate struct {
		Code string
	}
	DefectTag struct {
		Code string `validate:"aaaaaaaaaaaa"`
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
			validate, err := Validate(tt.in)
			assert.True(t, err == nil)
			assert.Equal(t, tt.countError, len(validate))
		})
	}
}

func TestValidateTypeErrorAndNameField(t *testing.T) {
	tests := []struct {
		in    interface{}
		err   error
		field string
	}{
		{App{Version: "123"}, ErrInvalidLength, "Version"},
		{Response{Code: 100, Body: ""}, ErrNotInSet, "Code"},
		{User{ID: "12345678901234567890123456789012345", Name: "", Age: 20, Email: "test@test.ru", Role: "admin", Phones: []string{"12345678901"}, meta: nil},
			ErrInvalidLength, "ID"},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 5, Email: "test@test.ru", Role: "admin", Phones: []string{"12345678901"}, meta: nil},
			ErrLess, "Age"},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 100, Email: "test@test.ru", Role: "admin", Phones: []string{"12345678901"}, meta: nil},
			ErrMax, "Age"},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 20, Email: "test", Role: "admin", Phones: []string{"12345678901"}, meta: nil},
			ErrNotMatchRegexp, "Email"},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 20, Email: "test@test.ru", Role: "", Phones: []string{"12345678901"}, meta: nil},
			ErrNotInSet, "Role"},
		{User{ID: "123456789012345678901234567890123456", Name: "", Age: 20, Email: "test@test.ru", Role: "admin", Phones: []string{"1234567890"}, meta: nil},
			ErrInvalidLength, "Phones"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			validate, err := Validate(tt.in)
			assert.True(t, err == nil)
			assert.Equal(t, 1, len(validate))
			assert.True(t, errors.Is(validate[0].Err, tt.err))
			assert.Equal(t, tt.field, validate[0].Field)
		})
	}
}

func TestDefectMetaInf(t *testing.T) {
	tests := []struct {
		in interface{}
	}{
		{DefectIn{Code: 1}},
		{DefectLen{Code: "1"}},
		{DefectTag{Code: "1"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			validate, err := Validate(tt.in)
			assert.True(t, err != nil)
			assert.True(t, validate == nil)
		})
	}
}

func TestIsStruct(t *testing.T) {
	validate, err := Validate("test")
	assert.Equal(t, 0, len(validate))
	assert.True(t, errors.Is(err, ErrNotStruct))
}

func TestEmpty(t *testing.T) {
	validate, err := Validate(EmptyValidate{Code: ""})
	assert.Equal(t, 0, len(validate))
	assert.True(t, err == nil)
}
