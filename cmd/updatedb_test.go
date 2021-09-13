package cmd

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestUpdatedb(t *testing.T) {
	com := Command{Output: "/var/lib/mlocate"}
	expected := fmt.Sprintf("%v",
		exec.Command(
			"updatedb",
			"-U",
			"/usr/bin",
			"--output",
			"/var/lib/mlocate/usr_bin.db"))
	actual := fmt.Sprintf("%v", com.Updatedb("/usr/bin"))
	if expected != actual {
		t.Fatalf("\ngot:  %v\nwant: %v", actual, expected)
	}
}

func TestReplaceUnder(t *testing.T) {
	expected := "/usr_bin"
	actual := replaceUnder("/usr/bin")
	if expected != actual {
		t.Fatalf("got: %v, want: %v", actual, expected)
	}
}
