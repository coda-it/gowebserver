package logger

import (
	"testing"
	"log"
	"os"
)

func TestSetup(t *testing.T) {
	t.Run("Should setup logger", func(t *testing.T) {
		expectedLogger := log.New(os.Stdout, "prefix: ", log.LstdFlags)
		expectedPrefix := expectedLogger.Prefix()

		Init("prefix")
		prefix := logger.Prefix()

		if expectedPrefix != prefix {
			t.Errorf("Prefixes don't match")
		}
	})
}