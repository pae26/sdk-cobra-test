package cmd

import (
	"testing"
)

func TestEditItem(t *testing.T) {
	t.Run("dry-run", func(t *testing.T) {
		err := editItem(false, "test-pae1", "test", "testtoken.yaml")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("unspecified title", func(t *testing.T) {
		expected := "ERROR: Set title of item with [-t] option."
		err := editItem(false, "test-pae1", "", "testtoken.yaml")
		if err.Error() != expected {
			t.Error(err)
		}
	})

	t.Run("apply", func(t *testing.T) {
		err := editItem(true, "test-pae1", "test3", "../testtoken.yaml")
		if err != nil {
			t.Error(err)
		}
	})
}
