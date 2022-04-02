package anchors

import "testing"

func TestEmptyDocumentMustReturnEmptyResponse(t *testing.T) {
	validEntriesToWalk := getAnchorsFromText("")
	if len(validEntriesToWalk) != 0 {
		t.Errorf("empty document must return empty entries")
	}
}

func TestSingleAnchor(t *testing.T) {
	validEntriesToWalk := getAnchorsFromText(`<a href="local.html">`)
	if len(validEntriesToWalk) != 1 {
		t.Errorf("must return one anchor element")
	}
}

func TestNestedAnchor(t *testing.T) {
	t.Run("single anchor", func(t *testing.T) {
		validEntriesToWalk := getAnchorsFromText(`<a href="local.html" <a hred="another">`)
		if len(validEntriesToWalk) != 1 {
			t.Errorf("must return one anchor element")
		}
	})

	t.Run("two anchors", func(t *testing.T) {
		doc := `<a href="local.html" <a hred="another.html"> text <a href="one.html" <a hred="two.html">`
		validEntriesToWalk := getAnchorsFromText(doc)
		if len(validEntriesToWalk) != 2 {
			t.Errorf("must return two anchor element")
		}
	})

	t.Run("multiple anchors", func(t *testing.T) {
		doc := `<a href="local.html" <a hred="another.html"> text <a href="one.html" <a hred="two.html"> other text <a href="four.html" <a hred="five.html">`
		validEntriesToWalk := getAnchorsFromText(doc)
		if len(validEntriesToWalk) != 3 {
			t.Errorf("must return 3 anchor element")
		}
	})
}

func TestGetAnchorTag(t *testing.T) {
	t.Run("valid anchor", func(t *testing.T) {
		validAnchor := getValidAnchorFromText(`<a href="local.html">`)
		if validAnchor == "" {
			t.Errorf("must return valid anchor element")
		}
	})

	t.Run("valid anchor nested tags", func(t *testing.T) {
		validAnchor := getValidAnchorFromText(`<a href="local.html" <a href="remote.html">`)
		if validAnchor != `<a href="local.html" ` {
			t.Errorf("must return valid anchor element")
		}
	})
}

func TestExtractHrefValidAnchorTag(t *testing.T) {
	t.Run("valid anchor", func(t *testing.T) {
		hrefText := getHrefAnchorTag(`<a href="local.html">`)
		if hrefText != "local.html" {
			t.Errorf("must return `local.html`")
		}
	})

	t.Run("valid href with nested anchors", func(t *testing.T) {
		hrefText := getHrefAnchorTag(`<a href="local.html" <a href="notlocal.html">`)
		if hrefText != "local.html" {
			t.Errorf("must return `local.html`")
		}
	})
}

func TestExtractHrefsToWalk(t *testing.T) {
	t.Run("one line doc", func(t *testing.T) {
		doc := `<a href="this is a document" <a href="this is an anchor nested">`

		href := GetWalkValidPages(doc)
		if href[0] != "this is a document" {
			t.Errorf("couldn't find the correct href")
		}
	})

	t.Run("complex doc", func(t *testing.T) {
		doc := `<h1>New section</h1>
<ul>
	<li><a href="first.html" <a href="nested_first.html"></li>
	<li><a href="second.html" <a href="nested_second.html" <a href="nested_nested_second.html"></li>
	<li><a href="third.html"></li>
</ul>`

		hrefs := GetWalkValidPages(doc)

		t.Run("first index", func(t *testing.T) {
			if hrefs[0] != "first.html" {
				t.Errorf("expected result for first index is wrong")
			}
		})

		t.Run("second index", func(t *testing.T) {
			if hrefs[1] != "second.html" {
				t.Errorf("expected result for second index is wrong")
			}
		})

		t.Run("third index", func(t *testing.T) {
			if hrefs[2] != "third.html" {
				t.Errorf("expected result for third index is wrong")
			}
		})
	})
}