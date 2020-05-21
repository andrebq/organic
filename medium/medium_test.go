package medium

import (
	"context"
	"testing"

	"github.com/andrebq/organic/cell"
	"github.com/andrebq/organic/internal/testutil"
)

// TestTypeCheck just validates if Medium matches the expected interface from cell
func TestTypeCheck(t *testing.T) {
	agar, err := NewAgar(context.Background(), "localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer agar.Close()
	c := cell.Grow("dummy", agar)
	_ = c
}

// TestMembraneIsValid checks if a membrane created by a medium is valid
func TestMembraneIsValid(t *testing.T) {
	agar, err := NewAgar(context.Background(), "localhost:6379")
	if err != nil {
		t.Fatal(err)
	}
	defer agar.Close()
	testutil.ExchangeMessage(t, "test_membrane_is_valid", agar)
}
