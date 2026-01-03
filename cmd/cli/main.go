package main

import (
	"os"

	"github.com/ChangChunCheng/webprism/internal/adapters/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
