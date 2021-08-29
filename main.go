// Run
// $ go test -v
package main

import (
	"fmt"
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

var opts = []Opt{
	{"./test/bin.db", "fstab"},
	{"./test/var.db", "fstab"},
	{"./test/etc.db", "fstab"},
	{"./test/usr.db", "fstab"},
}

func receiver(ch <-chan string) {
	for {
		s, ok := <-ch
		if ok == false {
			break
		}
		fmt.Println(s)
	}
}

func main() {
	var wg sync.WaitGroup // カウンタを宣言
	c := make(chan string)
	defer close(c) // main関数終了時にチャネル終了

	go receiver(c)
	for _, o := range opts {
		wg.Add(1) // カウンタの初期化
		go func(o Opt) {
			b, _ := exec.Command("locate", "-i", "-d", o.Dir, "--regex", o.Word).Output()
			out := strings.Split(string(b), "\n")
			for _, o := range out {
				if o == "" {
					break
				}
				c <- o
				time.Sleep(1 * time.Millisecond)
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
	for _, o := range opts {
		b, _ := exec.Command("locate", "-i", "-d", o.Dir, o.Word).Output()
		out := strings.Split(string(b), "\n")
		for _, o := range out {
			fmt.Println(o)
		}
	}
}
