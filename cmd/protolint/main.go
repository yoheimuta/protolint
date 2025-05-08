package main

import (
	"os"

	"github.com/yoheimuta/protolint/internal/cmd"
)

func main() {
	// Initialize the lint runner
	cmd.Initialize()

	os.Exit(int(
		cmd.Do(
			os.Args[1:],
			os.Stdout,
			os.Stderr,
		),
	))
}
