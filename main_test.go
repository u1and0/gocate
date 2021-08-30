package main

import (
	cmd "gocate/cmd"
	"sync"
	"testing"
)

func BenchmarkNormalLocate(b *testing.B) {
	args := []string{
		"-i",
		"-d",
		"./test/var.db:./test/etc.db:./test/usr.db",
		"--regex",
		".*pacman.*proto",
	}
	for i := 0; i < b.N; i++ {
		normalLocate(args)
	}
}

func BenchmarkParallelLocate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}

func TestMain(t *testing.T) {
	com = cmd.Command{
		Exe:  "/usr/sbin/locate",
		Args: []string{"fstab"},
		Wg:   sync.WaitGroup{},
	}
	c := make(chan string)

	go cmd.Receiver(c)
	dd := []string{
		"test/etc.db",
		"test/var.db",
		"test/usr.db",
	}
	for _, d := range dd {
		com.Wg.Add(1)
		com.Dir = d
		go com.Exec(c)
	}
	com.Wg.Wait()
}
