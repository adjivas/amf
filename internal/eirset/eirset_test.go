package util_test

import (
	"testing"
	"errors"

	eirset "github.com/free5gc/amf/internal/eirset"
	"github.com/stretchr/testify/assert"
)

func TestEirSet_WithFirstEir(t *testing.T) {
	set := eirset.New()
	expected := "http://127.0.0.8:8000"

	err_add := set.Add("http://127.0.0.8:8000")
	first, err := set.Next()

	assert.Equal(t, expected, first)
	assert.Nil(t, err_add)
	assert.Nil(t, err)

	// Same Eir
	expected_err_add := errors.New("EIR value already exists")
	err_add = set.Add("http://127.0.0.8:8000")
	second, err := set.Next()

	assert.Equal(t, expected, second)
	assert.Equal(t, expected_err_add, err_add)
	assert.Nil(t, err)
}

func TestEirSet_WithRemoveEir(t *testing.T) {
	set := eirset.New()
	expected := "http://127.0.0.8:8000"

	err_add := set.Add("http://127.0.0.8:8000")
	first, err := set.Next()

	assert.Equal(t, expected, first)
	assert.Nil(t, err_add)
	assert.Nil(t, err)

	// Remove Eir
	expected_next := ""
	expected_err_next := errors.New("EIR set is empty")

	err_remove := set.Remove("http://127.0.0.8:8000")
	next, err_next := set.Next()

	assert.Nil(t, err_remove)
	assert.Equal(t, expected_next, next)
	assert.Equal(t, expected_err_next, err_next)

	// Remove the same
	expected_err_remove := errors.New("EIR missing value")
	expected_err_next = errors.New("EIR set is empty")
	expected_next = ""

	err_remove = set.Remove("http://127.0.0.8:8000")
	next, err_next = set.Next()

	assert.Equal(t, expected_err_remove, err_remove)
	assert.Equal(t, expected_next, next)
	assert.Equal(t, expected_err_next, err_next)
}

func TestEirSet_WithEmptyEir(t *testing.T) {
	set := eirset.New()
	expected_next := ""
	expected_err_next := errors.New("EIR set is empty")

	first, err := set.Next()

	assert.Equal(t, expected_next, first)
	assert.Equal(t, expected_err_next, err)
}
