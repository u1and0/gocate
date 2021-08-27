package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	b, _ := exec.Command("locate", "-i", "-d", "./test/usr.db", "bin").Output()
	out := strings.Split(string(b), "\n")
	for _, o := range out {
		fmt.Printf("%s\n", o)
	}
	fmt.Printf("Length %d", len(out))
}
