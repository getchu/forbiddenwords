package lib

import (
	"forbiddenwords/lib"
	"testing"
)

//测试Set性能
func BenchmarkCache_Set(b *testing.B) {
	c := lib.NewCache()
	c.Set("k1", "v1")
	for i := 0; i < b.N; i++ {
		c.Set("k1", "v1")
	}
}

//测试Get性能
func BenchmarkCache_Get(b *testing.B) {
	c := lib.NewCache()
	c.Set("k1", "v1")
	for i := 0; i < b.N; i++ {
		c.Get("k1")
	}
}
