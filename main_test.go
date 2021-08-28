package main

import (
	"testing"
)

func BenchmarkNormalLocate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		locateBin0()
		locateUsr0()
	}
}

func BenchmarkParallelLocate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}
