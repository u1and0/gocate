package cmd

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path"
	"strings"
)

// FileTree : ioutil.ReadDir() root + directories set
type FileTree struct {
	Pairent string
	Dirs    []fs.FileInfo
}

// Updatedb generate command updatedb [OPTION]...
func (com *Command) Updatedb(f string) *exec.Cmd {
	opt := []string{"-U", f, "--output", path.Join(com.Output, replaceUnder(f)) + ".db"}
	opt = append(opt, com.Args...) // `updatedb` command other option
	return exec.Command("updatedb", opt...)
}

func replaceUnder(s string) string { // <= /usr/bin
	return "/" + strings.ReplaceAll(s[1:], "/", "_") // => /usr_bin
}

// DirectoryFilter : return directories array. if get not directory file, print warning.
func (ft *FileTree) DirectoryFilter(dd []string) []string {
	for _, d := range ft.Dirs { // <= FileTree{Pairent:/usr, Dirs: []fs.FileInfo{bin, lib, ...}}
		concatdir := path.Join(ft.Pairent, d.Name())
		if d.IsDir() {
			dd = append(dd, concatdir) // => dd = /usr/bin /usr/lib ... /etc/iptables
		} else {
			fmt.Printf("warning: %s is not directory, it will be ignored for indexing.\n", concatdir)
		}
	}
	return dd
}
