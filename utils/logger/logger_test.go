package logger

import (
	"testing"
	"log"
	"os"
)

func TestLog(t *testing.T) {
	t.Run("Should setup logger", func(t *testing.T) {
		expectedLogger := log.New(os.Stdout, buildPrefix("prefix"), log.LstdFlags)
		Setup("prefix")

		if expectedLogger != logger {
			t.Errorf("", )
		}
	})
}
