package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bep/semverpair"
	"github.com/spf13/pflag"
	"golang.org/x/mod/semver"
)

func toVersion(s string) (semverpair.Version, error) {
	// convert to canonical semver string (results in an empty string if invalid)
	s = semver.Canonical(s)

	// error immediately on blank string
	if s == "" {
		return semverpair.Version{}, fmt.Errorf("\"%s\" was not a valid semver string", s)
	}

	// trim any build and prerelease suffix
	s = strings.TrimSuffix(s, semver.Build(s))
	s = strings.TrimSuffix(s, semver.Prerelease(s))

	// strip leading "v" and split into parts
	split := strings.Split(strings.TrimPrefix(s, "v"), ".")

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

func shift(a []string) []string {
	if len(a) == 0 {
		return a
	}

	if len(a) == 1 {
		return make([]string, 0)
	}

	return a[1:]
}

func parsePositionalArgs(args []string, first, second string) (string, string) {
	if first == "" && len(args) >= 1 {
		first = args[0]
		args = shift(args)
	}

	if second == "" && len(args) >= 1 {
		second = args[0]
	}

	return first, second
}

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

func main() {
	var first, second, decode string

	pflag.StringVarP(&decode, "decode", "d", "", "Decode a semver triplet into the original version pair")
	pflag.StringVarP(&first, "first", "f", "", "First version to encode (also can be provided as a positional argument)")
	pflag.StringVarP(&second, "second", "s", "", "Second version to encode (also can be provided as a positional argument)")
	pflag.Parse()

	// parse any positional args
	first, second = parsePositionalArgs(pflag.Args(), first, second)

	if decode != "" && (first != "" || second != "") {
		fmt.Fprintf(os.Stderr, "Error: decode, first and second options are mutually exclusive\n")
		os.Exit(1)
	}

	if decode != "" {
		pair, err := decodePair(decode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("First = %s; Second = %s\n", pair.First, pair.Second)
	} else {
		pair, err := encodePair(first, second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", semverpair.Encode(pair))
	}
}
