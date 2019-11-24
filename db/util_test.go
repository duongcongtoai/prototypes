package db

import (
	"testing"
)

func TestOpenDb(t *testing.T) {
	_, _, err := OpenDb()
	if err != nil {
		t.Errorf("error return after calling the function: %v", err)
	}
}
