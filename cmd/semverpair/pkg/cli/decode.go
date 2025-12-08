package cli

import (
	"context"
	"fmt"

	"github.com/andrewheberle/simplecommand"
	"github.com/bep/simplecobra"
)

type decodeCommand struct {
	version string

	*simplecommand.Command
}

func (c *decodeCommand) Init(cd *simplecobra.Commandeer) error {
	if err := c.Command.Init(cd); err != nil {
		return err
	}

	cmd := cd.CobraCommand
	cmd.Flags().StringVarP(&c.version, "version", "v", "", "Version to decode")

	return nil
}

func (c *decodeCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	if err := c.Command.PreRun(this, runner); err != nil {
		return err
	}

	args := parsePositionalArgs(this.CobraCommand.Flags().Args(), c.version)
	if args[0] == "" {
		return fmt.Errorf("version string is required")
	}
	c.version = args[0]

	return nil
}

func (c *decodeCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	// decode version
	pair, err := decodePair(c.version)
	if err != nil {
		return fmt.Errorf("could not decode: %w", err)
	}

	fmt.Printf("First = %s; Second = %s\n", pair.First, pair.Second)

	return nil
}
