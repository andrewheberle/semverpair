package cli

import (
	"context"
	"fmt"

	"github.com/andrewheberle/simplecommand"
	"github.com/bep/semverpair"
	"github.com/bep/simplecobra"
)

type encodeCommand struct {
	first  string
	second string

	*simplecommand.Command
}

func (c *encodeCommand) Init(cd *simplecobra.Commandeer) error {
	if err := c.Command.Init(cd); err != nil {
		return err
	}

	cmd := cd.CobraCommand
	cmd.Flags().StringVarP(&c.first, "first", "f", "", "First version to encode (also can be provided as a positional argument)")
	cmd.Flags().StringVarP(&c.second, "second", "s", "", "Second version to encode (also can be provided as a positional argument)")

	return nil
}

func (c *encodeCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	if err := c.Command.PreRun(this, runner); err != nil {
		return err
	}

	args := parsePositionalArgs(this.CobraCommand.Flags().Args(), c.first, c.second)
	if args[0] == "" || args[1] == "" {
		return fmt.Errorf("first and second version string are required")
	}
	c.first = args[0]
	c.second = args[1]

	return nil
}

func (c *encodeCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	// decode version
	pair, err := encodePair(c.first, c.second)
	if err != nil {
		return fmt.Errorf("could not encode: %w", err)
	}

	fmt.Printf("%s\n", semverpair.Encode(pair))

	return nil
}
