package main

import (
	"testing"
	"time"
)

func TestNormalLocate(t *testing.T) {
	normalLocate()
}

func TestParallelLocate(t *testing.T) {
	parallelLocate()
	time.Sleep(9 * time.Millisecond)
}
