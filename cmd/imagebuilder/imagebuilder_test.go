package main

import (
	"os"
	"testing"
)

//TODO(illia) test for commandline utility
func TestBuild(t *testing.T) {
	args := []string{"build", "ccache"}
	os.Args = append(os.Args, args...)
	MockExit(func() {

	})
	//main()
}

func MockExit(checkOnExit func()) {
	exit = func(status int) {
		checkOnExit()
	}
}
