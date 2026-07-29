package gowebserver

import (
	"testing"
	"time"
)

func TestResolveTimeout(t *testing.T) {
	def := 30 * time.Second

	if got := resolveTimeout(0, def); got != def {
		t.Errorf("zero should use default: got %v, want %v", got, def)
	}
	if got := resolveTimeout(5*time.Second, def); got != 5*time.Second {
		t.Errorf("positive should be used as-is: got %v", got)
	}
	if got := resolveTimeout(-1, def); got != 0 {
		t.Errorf("negative should disable the timeout (0): got %v", got)
	}
}
