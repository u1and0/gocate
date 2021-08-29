// Run
// $ go test -v
package main

import (
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Opt : locate command options
type Opt struct {
	Dir  string
	Word string
}

const (
	opts string = "./test/bin.db:./test/var.db:./test/etc.db:./test/usr.db"
	word string = "fstab"
)

func receiver(ch <-chan string) {
	for {
		_, ok := <-ch
		if ok == false {
			break
		}
		// fmt.Println(s)
	}
}

func main() {
	var wg sync.WaitGroup // カウンタを宣言
	c := make(chan string)
	defer close(c) // main関数終了時にチャネル終了

	go receiver(c)
	for _, o := range strings.Split(opts, ":") {
		wg.Add(1) // カウンタの初期化
		go func(o string) {
			b, _ := exec.Command("locate", "-i", "-d", o, "--regex", word).Output()
			out := strings.Split(string(b), "\n")
			for _, s := range out {
				if s == "" {
					break
				}
				c <- s
			}
			wg.Done() // カウンタを減算
		}(o)
		wg.Wait() // カウンタが0になるまでブロック
	}
}

// // Locate excutes locate command
// func Locate(o Opt, wg *sync.WaitGroup) {
// 	b, _ := exec.Command("locate", "-i", "-d", o.Dir, o.Word).Output()
// 	out := strings.Split(string(b), "\n")
// 	for _, o := range out {
// 		time.Sleep(1 * time.Microsecond)
// 		fmt.Println(o)
// 	}
// 	wg.Done() // カウンタを減算
// }

// Nomral locate command for benchmark
func normalLocate() {
	b, _ := exec.Command("locate", "-i", "-d", opts, word).Output()
	out := strings.Split(string(b), "\n")
	for range out {
		// fmt.Println(o)
		time.Sleep(1 * time.Microsecond)
	}
}

/*
$ go test -bench .
goos: linux
goarch: amd64
pkg: speedtest/src/github.com/u1and0/gocate
cpu: Intel(R) Core(TM) i5-8400 CPU @ 2.80GHz
BenchmarkNormalLocate-6               44          28034939 ns/op
BenchmarkParallelLocate-6             20          67518023 ns/op
PASS
ok      speedtest/src/github.com/u1and0/gocate  2.683s

普通のlocateの方が3倍早い
*/
