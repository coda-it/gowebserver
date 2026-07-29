package url

import (
	"regexp"
	"testing"
)

func TestPatternToRegExp(t *testing.T) {
	t.Run("Should translate route pattern into URL", func(t *testing.T) {
		inputURL := "/path1/path2/{id}/"
		finalURL := PatternToRegExp(inputURL)

		if finalURL != `^\/path1\/path2(\/([0-9a-zA-Z_-])*)?\/$` {
			t.Errorf("Transformed url is not correct, got: %s.", finalURL)
		}
	})

	t.Run("Should match a hyphenated slug in a parameter segment", func(t *testing.T) {
		pattern := regexp.MustCompile(PatternToRegExp("/post/{id}"))

		if !pattern.MatchString("/post/my-first-post-2") {
			t.Error("Expected pattern to match a hyphenated slug")
		}
	})

	t.Run("Should match an object id in a parameter segment", func(t *testing.T) {
		pattern := regexp.MustCompile(PatternToRegExp("/post/{id}"))

		if !pattern.MatchString("/post/64b8f0c2a1e4d5f6a7b8c9d0") {
			t.Error("Expected pattern to match an object id")
		}
	})

	t.Run("Should match the bare path when the parameter segment is omitted", func(t *testing.T) {
		pattern := regexp.MustCompile(PatternToRegExp("/post/{id}"))

		if !pattern.MatchString("/post") {
			t.Error("Expected pattern to match the bare path without a parameter")
		}
	})

	t.Run("Should not match a deeper unrelated path", func(t *testing.T) {
		pattern := regexp.MustCompile(PatternToRegExp("/post/{id}"))

		if pattern.MatchString("/post/my-post/extra") {
			t.Error("Expected pattern not to match a deeper path")
		}
	})
}
