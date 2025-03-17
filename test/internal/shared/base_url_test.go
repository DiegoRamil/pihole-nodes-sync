package shared

import (
	"testing"

	"github.com/DiegoRamil/pihole-nodes-sync/internal/shared"
)

func TestConcatBaseUrlAndUri(t *testing.T) {
	baseUrl := ""
	uri := "test"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but function did not panic")
		}
	}()
	shared.ConcatBaseUrlAndUri(baseUrl, uri)
}

func TestConcatBaseUrlAndUriWithNormalValues(t *testing.T) {
	baseUrl := "http://localhost:332"
	uri := "/api/test"
	res := shared.ConcatBaseUrlAndUri(baseUrl, uri)
	expected := "http://localhost:332/api/test"

	if res != expected {
		t.Errorf("expected %s but got %s", expected, res)
	}
}

func TestConcatBaseUrlAndUriWithSlashInBaseUrl(t *testing.T) {
	baseUrl := "http://localhost:332/"
	uri := "/api/test"
	res := shared.ConcatBaseUrlAndUri(baseUrl, uri)
	expected := "http://localhost:332/api/test"

	if res != expected {
		t.Errorf("expected %s but got %s", expected, res)
	}
}
