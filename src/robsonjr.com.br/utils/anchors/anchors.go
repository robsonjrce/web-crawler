package anchors

import (
	"regexp"
)

// getAnchorsFromText: will search on the text for anchor <a> tags
func getAnchorsFromText(documentText string) []string {
	regexAnchors := regexp.MustCompile(`<a(?:[^>].*?)>`)
	validAnchors := regexAnchors.FindAllString(documentText, -1)
	return validAnchors
}

// getValidAnchorFromText: will extract valid anchor <a> tag from text
func getValidAnchorFromText(tagText string) string {
	// we are matching the smallest possible valid anchor, that means
	// nested `<a` are disconsidered
	regexValidAnchor := regexp.MustCompile(`<a.*?[<>]`)
	validAnchor := regexValidAnchor.Find([]byte(tagText))
	validAnchorWithoutLastChar := validAnchor[:len(validAnchor)-1]
	return string(validAnchorWithoutLastChar)
}

// getHrefAnchorTag: will extract href attribute from valid anchor <a> tag
func getHrefAnchorTag(anchorText string) string {
	regexHrefAttr := regexp.MustCompile(`<a.*?href=['"]([^"]*?)['"].*`)
	textHref := regexHrefAttr.FindSubmatch([]byte(anchorText))
	if textHref == nil || len(textHref) != 2 {
		return ""
	}
	return string(textHref[1])
}

// GetWalkValidPages: will return href attributes for all valid anchor <a> tags
func GetWalkValidPages(document string) []string {
	pagesToWalk := []string{}

	anchors := getAnchorsFromText(document)
	for _, anchor := range anchors {
		sanitizeAnchor := getValidAnchorFromText(anchor)
		href := getHrefAnchorTag(sanitizeAnchor)
		if href != "" {
			pagesToWalk = append(pagesToWalk, href)
		}
	}

	return pagesToWalk
}