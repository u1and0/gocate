// Run
// $ go test -v
package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// nothing to do
}

func locateBin() {
	b, _ := exec.Command("locate", "-i", "-d", "./test/bin.db", "ls").Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%s\n", o)
	}
}

func locateUsr() {
	b, _ := exec.Command("locate", "-i", "-d", "./test/usr.db", "ing.hpp").Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%s\n", o)
	}
}

func normalLocate() {
	locateUsr()
	locateBin()
}

func parallelLocate() {
	go locateUsr()
	go locateBin()
}
