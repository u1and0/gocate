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
	time.Sleep(4 * time.Millisecond)
}
