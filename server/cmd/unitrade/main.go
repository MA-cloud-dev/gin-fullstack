package main

import (
	"fmt"
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/unitrade"
)

func main() {
	app, err := unitrade.NewApp(os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize unitrade: %v\n", err)
		os.Exit(1)
	}
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "unitrade: %v\n", err)
		os.Exit(1)
	}
}
