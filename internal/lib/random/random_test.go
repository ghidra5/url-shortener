package random_test

import (
	"testing"
	"url-shortener/internal/lib/random"
)

// TODO: add more tests
func TestNewRandomString(t *testing.T) {
	length := 6
	randomString := random.NewRandomString(length)

	if len(randomString) != length {
		t.Errorf("expected length %d, got %d", length, len(randomString))
	}

}
