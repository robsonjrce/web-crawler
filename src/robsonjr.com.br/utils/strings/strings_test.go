package strings

import "testing"

func TestEmptyDocumentMustReturnEmptyResponse(t *testing.T) {
	validEntriesToWalk := GetValidAnchorHrefFromText("")
	if len(validEntriesToWalk) != 0 {
		t.Errorf("empty document must return empty entries")
	}
}