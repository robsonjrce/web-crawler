package validation

import "testing"

func TestIsChildrenUrl(t *testing.T) {
	t.Run("default schema", func(t *testing.T) {
		if !IsChildrenUrl("http://google.com", "google.com") {
			t.Errorf("the url without schema should match http")
		}
	})

	t.Run("different schemas", func(t *testing.T) {
		if !IsChildrenUrl("https://google.com", "google.com") {
			t.Errorf("the schema must not influentiate on children check")
		}
	})

	t.Run("children 1", func(t *testing.T) {
		if !IsChildrenUrl("http://google.com", "google.com/chrome") {
			t.Errorf("the url without schema should match http")
		}
	})

	t.Run("children 2", func(t *testing.T) {
		if !IsChildrenUrl("http://google.com?date=20220401", "google.com/chrome") {
			t.Errorf("the url without schema should match http")
		}
	})
}
