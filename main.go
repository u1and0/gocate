// Run
// $ go test -v
package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	// nothing to do
}

func locateBin() {
	b, _ := exec.Command("locate", "-i", "-d", "./test/bin.db", "ls").Output()
	time.Sleep(1 * time.Microsecond)
	fmt.Printf("%v", string(b))
}

func locateUsr() {
	b, _ := exec.Command("locate", "-i", "-d", "./test/usr.db", "ing.hpp").Output()
	time.Sleep(1 * time.Microsecond)
	fmt.Printf("%v", string(b))
}

func normalLocate() {
	locateUsr()
	locateBin()
}

func parallelLocate() {
	go locateUsr()
	go locateBin()
}
