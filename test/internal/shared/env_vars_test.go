package shared

import (
	"testing"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/shared"
)

func TestPanicOnNullValue(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but function did not panic")
		}
	}()

	shared.RetrieveEnvVar("TEST")
}

func TestRetrieveEnvVar(t *testing.T) {
	t.Setenv("TEST", "test")
	v := shared.RetrieveEnvVar("TEST")
	if v != "test" {
		t.Errorf("Expected test, got %s", v)
	}
}
