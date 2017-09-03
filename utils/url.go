package utils

import (
	"regexp"
	"strings"
)

func UrlPatternToRegExp(urlPattern string) string {
	mapParameters := func(urlItem string) string {
		return "(/([0-9a-zA-Z])*)?"
	}
	wrapUrl := func(url string) string {
		return `^` + url + `$`
	}

	paramsRegExp, _ := regexp.Compile(`/{[a-zA-Z0-9]*}`)

	finalUrl := paramsRegExp.ReplaceAllStringFunc(urlPattern, mapParameters)
	finalUrl = strings.Replace(finalUrl, "/", `\/`, -1)
	finalUrl = wrapUrl(finalUrl)

	return finalUrl
}
