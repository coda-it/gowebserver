package url

import "testing"

func TestPatternToRegExp(t *testing.T) {
	t.Run("Should translate route pattern into URL", func(t *testing.T) {
		inputURL := "/path1/path2/{id}/"
		expectedRegExp := "/path1/path2/{[.]*}/"
		finalURL := PatternToRegExp(inputURL)

		if finalURL != `^\/path1\/path2(\/([0-9a-zA-Z])*)?\/$` {
			t.Errorf("Transformed url is not correct, got: %s, want: %s.",
				finalURL, expectedRegExp)
		}
	})
}
