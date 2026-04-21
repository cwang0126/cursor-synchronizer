package main

import (
	"fmt"
	"os"

	"github.com/cwang0126/cursor-synchronizer/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
