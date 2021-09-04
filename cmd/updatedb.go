package cmd

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path"
	"strings"
)

// Updatedb generate command updatedb [OPTION]...
func (com *Command) Updatedb(f string) *exec.Cmd {
	opt := []string{
		"-U",
		fmt.Sprintf("%s", f),
		"--output",
		path.Join(com.Gocatedbpath, replaceAnder(f)) + ".db",
		// fmt.Sprintf("%s%s.db", com.Gocatedbpath, replaceAnder(f)),
		// `pwd`/.gocate is temp directory
	}
	return exec.Command("updatedb", opt...)
}

func replaceAnder(s string) string {
	return "/" + strings.ReplaceAll(s[1:], "/", "_")
}

// FileTree : ioutil.ReadDir() root + directories set
type FileTree struct {
	Root string
	Dirs []fs.FileInfo
}
