package lib

import (
	"forbiddenwords/lib"
	"testing"
)

func TestNewCache(t *testing.T) {
	c := lib.NewCache()
	c.Set("k1", "v1")
	value, err := c.Get("k1")
	if err != nil {
		t.Fatalf("get error: %s", err.Error())
	}
	if value != "v1" {
		t.Fatalf("get value invalid: %s / %s", value, "v1")
	}
}
