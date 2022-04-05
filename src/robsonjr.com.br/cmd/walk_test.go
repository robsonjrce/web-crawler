package cmd

import "testing"

func TestEnqueue(t *testing.T) {
	pendingUrls := make([]string, 0)
	workingUrls := make(map[string]bool)

	pendingUrls = append(pendingUrls, "https://go.dev/")
	pendingUrls = append(pendingUrls, "https://go.dev/")
	pendingUrls = append(pendingUrls, "https://go.dev/")
	pendingUrls = append(pendingUrls, "https://go.dev/")

	t.Run("first call", func(t *testing.T) {
		nextUrl := getNextUrlToWalk(&pendingUrls, &workingUrls)

		if nextUrl != "https://go.dev/" {
			t.Errorf("the nextUrl is not the value excpected")
		}
	})

	t.Run("remaining items", func(t *testing.T) {
		if len(pendingUrls) != 3 {
			t.Errorf("the remaining items should exist")
		}
	})

	t.Run("duplicated call", func(t *testing.T) {
		nextUrl := getNextUrlToWalk(&pendingUrls, &workingUrls)

		if nextUrl != "" {
			t.Errorf("the nextUrl should be empty as it was already visited: %v", nextUrl)
		}
	})

	t.Run("no more itens on stack call", func(t *testing.T) {
		if len(pendingUrls) != 0 {
			t.Errorf("the expected size for pending urls is 0")
		}
	})

}
