package main

import (
	"os"

	"github.com/ipfred/cc-hooks-notify/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
