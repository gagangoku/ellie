package main

import "os"

var FILE_FLAGS = os.O_CREATE | os.O_WRONLY | os.O_TRUNC

type CleanupFn func()

func noop(x ...any) {
	// Does nothing
}
