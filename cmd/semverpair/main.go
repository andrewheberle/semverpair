package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andrewheberle/semverpair/cmd/semverpair/pkg/cli"
)

func main() {
	if err := cli.Execute(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
