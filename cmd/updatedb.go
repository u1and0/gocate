package cmd

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
)

// Updatedb generate command updatedb [OPTION]...
func (c *Command) Updatedb(f fs.FileInfo) *exec.Cmd {
	if !f.IsDir() {
		// e := fmt.Sprintf("do not must be included file %s", f.Name())
		// return exec.Cmd{}, errors.New(e) //fmt.Errorf("do not must be included file %s", f.Name())
		fmt.Println("warning: ", f.Name(), "will be ignored.")
	}
	pairent := filepath.Dir(f.Name())
	opt := []string{
		"-U",
		fmt.Sprintf("%s/%s", pairent, f.Name()),
		"--output",
		fmt.Sprintf("%s/%s_%s.db", c.Gocatedbpath, pairent, f.Name()),
		// `pwd`/.gocate is temp directory
	}
	return exec.Command("updatedb", opt...)
	// fmt.Println(c)
	// command.Run()
}
