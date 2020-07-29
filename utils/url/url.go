package url

import (
	"regexp"
	"strings"
)

// PatternToRegExp - translates URL pattern to RegExp
func PatternToRegExp(urlPattern string) string {
	mapParameters := func(urlItem string) string {
		return "(/([0-9a-zA-Z])*)?"
	}
	wrapURL := func(url string) string {
		return `^` + url + `$`
	}

	paramsRegExp, _ := regexp.Compile(`/{[a-zA-Z0-9]*}`)

	finalURL := paramsRegExp.ReplaceAllStringFunc(urlPattern, mapParameters)
	finalURL = strings.Replace(finalURL, "/", `\/`, -1)
	finalURL = wrapURL(finalURL)

	return finalURL
}
