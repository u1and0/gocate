package main

import (
	"sync"
	"testing"

	cmd "github.com/u1and0/gocate/cmd"
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
	var (
		wg  = sync.WaitGroup{}
		com = cmd.Command{
			Args: []string{"--regex", "'lib.*id$'"},
		}
		ch = make(chan string)
	)

	go cmd.Receiver(ch)
	dd := []string{
		"test/etc.db",
		"test/var.db",
		"test/usr.db",
	}
	for _, d := range dd {
		wg.Add(1)
		go func(d string, ch chan string) {
			defer wg.Done()
			c := com.Locate(d)
			cmd.Run(*c, ch)
		}(d, ch)
	}
	wg.Wait()
}

func Test_arrayFieldDbPath(t *testing.T) {
	af := arrayField{"/usr", "/var"}
	expected := []string{ // $ ls -d /usr/* /var/*
		"/usr/bin",
		"/usr/include",
		"/usr/lib",
		"/usr/lib32",
		// "/usr/lib64", <- symbolic link
		"/usr/local",
		// "/usr/sbin", <- symbolic link
		"/usr/share",
		"/usr/src",
		"/var/cache",
		"/var/db",
		"/var/empty",
		"/var/games",
		"/var/lib",
		"/var/local",
		// "/var/lock", <- symbolic link
		"/var/log",
		// "/var/mail", <- symbolic link
		"/var/opt",
		// "/var/run", <- symbolic link
		"/var/spool",
		"/var/tmp",
	}
	actual := af.Dbpath()
	for i, e := range expected {
		if e != actual[i] {
			t.Fatalf("%s,%s,\ngot:  %v\nwant: %v", actual[i], e, actual, expected)
		}
	}
}
