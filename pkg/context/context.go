package context

import (
	"io"
	"log"
	"os"
	"os/exec"
)

var (
	// SystemCall function to do system call
	SystemCall = func(dir string, cmd string, args []string, out, err io.Writer) error {
		command := exec.Cmd{
			Dir:    dir,
			Path:   cmd,
			Stderr: err,
			Stdout: out,
			Args:   append([]string{cmd}, args...),
		}

		log.Printf("Command running: %s %s", cmd, args)

		return command.Run()
	}

	//TODO(illia)implement port assigning and releasing
	lastUsedPort = 5000
	GetFreePort  = func() int {
		lastUsedPort = lastUsedPort + 1
		return lastUsedPort
	}

	GetHostName = func() string {
		if name, err := os.Hostname(); err != nil {
			log.Printf("Couldn't get hostname err: %s", err)
			return ""
		} else {
			return name
		}
	}
)
