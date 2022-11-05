package main

import (
	"os"

	"github.com/dukeofdisaster/podbeat/cmd"

	_ "github.com/dukeofdisaster/podbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
