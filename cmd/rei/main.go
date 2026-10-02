package main

import (
	"os"

	"github.com/k1wi777/my-harness-SDD/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
