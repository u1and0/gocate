package cmd

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestLocate(t *testing.T) {
	com := Command{Args: []string{"keyword"}}
	expected := fmt.Sprintf("%v",
		exec.Command(
			"locate",
			"-d",
			"/var/lib/mlocate/test.db",
			"keyword",
		))
	actual := fmt.Sprintf("%v", com.Locate("/var/lib/mlocate/test.db"))
	if expected != actual {
		t.Fatalf("\ngot:  %v\nwant: %v", actual, expected)
	}
}
