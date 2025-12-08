package cli

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/andrewheberle/simplecommand"
	"github.com/bep/simplecobra"
)

type versionCommand struct {
	version string

	*simplecommand.Command
}

func (c *versionCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("error getting version")
	}

	c.version = bi.Main.Version

	return nil
}

func (c *versionCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	fmt.Printf("%s version: %s\n", cd.Root.CobraCommand.Name(), c.version)

	return nil
}
