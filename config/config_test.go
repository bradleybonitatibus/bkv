package config

import "testing"

func TestConfig(t *testing.T) {
	t.Run("Test that config can create the server configuration", func(t *testing.T) {
		c := NewConfig(8080, "0.0.0.0")
		if c.String() != "0.0.0.0:8080" {
			t.Error("Failed to format string correctly")
		}
	})
}
