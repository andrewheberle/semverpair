package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/andrewheberle/simplecommand"
	"github.com/bep/semverpair"
	"github.com/bep/simplecobra"
	"golang.org/x/mod/semver"
)

type rootCommand struct {
	*simplecommand.Command
}

func (c *rootCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	cd.CobraCommand.Usage()

	return fmt.Errorf("run one of the available sub-commands")
}

func Execute(ctx context.Context, args []string) error {
	root := &rootCommand{
		Command: simplecommand.New("semverpair", "Encode or decode semver strings"),
	}
	root.SubCommands = []simplecobra.Commander{
		&encodeCommand{
			Command: simplecommand.New("encode", "Encode a pair of semver strings"),
		},
		&decodeCommand{
			Command: simplecommand.New("decode", "Decode a semver string"),
		},
		&versionCommand{
			Command: simplecommand.New("version", "Show version"),
		},
	}

	// Set up simplecobra
	x, err := simplecobra.New(root)
	if err != nil {
		return err
	}

	// run command with the provided args
	if _, err := x.Execute(ctx, args); err != nil {
		return err
	}

	return nil
}

// Parses provided args and sets the values of to the provided values if not already set
func parsePositionalArgs(args []string, vars ...string) []string {
	result := make([]string, len(vars))
	for n, v := range vars {
		if v == "" && len(args) >= 1 {
			result[n] = args[0]
			args = shift(args)
		} else {
			result[n] = v
		}
	}

	return result
}

// Shifts a slice by removing the first element and returning the remining slice
func shift(a []string) []string {
	if len(a) == 0 {
		return a
	}

	if len(a) == 1 {
		return make([]string, 0)
	}

	return a[1:]
}

// Encodes a pair of valid semantic version strings and returns a semverpair.Pair
func encodePair(first, second string) (semverpair.Pair, error) {
	// check we have two versions to combine
	if first == "" || second == "" {
		return semverpair.Pair{}, fmt.Errorf("both first and second versions must be set")
	}

	// add v prefix if missing
	if !strings.HasPrefix(first, "v") {
		first = "v" + first
	}
	if !strings.HasPrefix(second, "v") {
		second = "v" + second
	}

	// decode version strings
	v1, err := toVersion(first)
	if err != nil {
		return semverpair.Pair{}, fmt.Errorf("error decoding first version: %w", err)
	}
	v2, err := toVersion(second)
	if err != nil {
		return semverpair.Pair{}, fmt.Errorf("error decoding second version: %w", err)
	}

	return semverpair.Pair{First: v1, Second: v2}, nil
}

// Decodes a semantic version string and returns a semverpair.Pair of the original versions
func decodePair(decode string) (semverpair.Pair, error) {
	// add v prefix if missing
	if !strings.HasPrefix(decode, "v") {
		decode = "v" + decode
	}

	v, err := toVersion(decode)
	if err != nil {
		return semverpair.Pair{}, fmt.Errorf("could not decode version: %w", err)
	}

	return semverpair.Decode(v), nil
}

// Converts a valid semantic version string to a semverpair.Version
func toVersion(s string) (semverpair.Version, error) {
	// convert to canonical semver string and error immediately on blank string
	v := semver.Canonical(s)
	if v == "" {
		return semverpair.Version{}, fmt.Errorf("\"%s\" was not a valid semver string", s)
	}

	// trim any build and prerelease suffix
	v = strings.TrimSuffix(v, semver.Build(s))
	v = strings.TrimSuffix(v, semver.Prerelease(s))

	// strip leading "v" and split into parts
	split := strings.Split(strings.TrimPrefix(v, "v"), ".")

	// convert to integers
	major, err := strconv.Atoi(split[0])
	if err != nil {
		// this should never happen
		panic(err)
	}
	minor, err := strconv.Atoi(split[1])
	if err != nil {
		// this should never happen
		panic(err)
	}
	patch, err := strconv.Atoi(split[2])
	if err != nil {
		// this should never happen
		panic(err)
	}

	return semverpair.Version{Major: major, Minor: minor, Patch: patch}, nil
}
