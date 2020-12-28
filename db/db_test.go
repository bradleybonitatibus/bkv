package db

import testing "testing"

func TestDatabase(t *testing.T) {
	var db Database
	t.Run("Test constructor of database", func(t *testing.T) {
		db = NewDatabase()
		if db == nil {
			t.Error("Failed to construct database")
		}
	})
	t.Run("Test GET returns error when key is not present", func(t *testing.T) {
		val, err := db.Get("Fake")
		if err == nil || val != "" {
			t.Error("Expected err to be non-nil")
		}
	})
	t.Run("Test GET after SET", func(t *testing.T) {
		k, v := "Hello", "World"
		db.Set(k, v)
		val, err := db.Get(k)
		if err != nil || val != v {
			t.Error("Expected error to not be nil and value to be set correctly")
		}
	})
}
