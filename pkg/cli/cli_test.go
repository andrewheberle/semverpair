package cli

import (
	"slices"
	"testing"

	"github.com/bep/semverpair"
)

func Test_parsePositionalArgs(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		first  string
		second string
		want   []string
	}{
		{"both set with no args", []string{}, "a", "b", []string{"a", "b"}},
		{"both set with one arg", []string{"other"}, "a", "b", []string{"a", "b"}},
		{"both set with two args", []string{"other", "again"}, "a", "b", []string{"a", "b"}},
		{"both  set with three args", []string{"other", "again", "ignored"}, "a", "b", []string{"a", "b"}},
		{"none set with no args", []string{}, "", "", []string{"", ""}},
		{"none set with one args", []string{"first"}, "", "", []string{"first", ""}},
		{"none set with two args", []string{"first", "second"}, "", "", []string{"first", "second"}},
		{"none set with three args", []string{"first", "second", "ignored"}, "", "", []string{"first", "second"}},
		{"first set with no args", []string{}, "firstset", "", []string{"firstset", ""}},
		{"first set with one args", []string{"second"}, "firstset", "", []string{"firstset", "second"}},
		{"first set with two args", []string{"second", "ignored"}, "firstset", "", []string{"firstset", "second"}},
		{"first set with three args", []string{"second", "ignored1", "ignored2"}, "firstset", "", []string{"firstset", "second"}},
		{"second set with no args", []string{}, "", "secondset", []string{"", "secondset"}},
		{"second set with one args", []string{"first"}, "", "secondset", []string{"first", "secondset"}},
		{"second set with two args", []string{"first", "ignored"}, "", "secondset", []string{"first", "secondset"}},
		{"second set with three args", []string{"first", "ignored1", "ignored2"}, "", "secondset", []string{"first", "secondset"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePositionalArgs(tt.args, tt.first, tt.second)
			if !slices.Equal(tt.want, got) {
				t.Errorf("parsePositionalArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shift(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		want []string
	}{
		{"nil slice", nil, nil},
		{"empty slice", make([]string, 0), make([]string, 0)},
		{"one element", []string{"a"}, make([]string, 0)},
		{"two elements", []string{"a", "b"}, []string{"b"}},
		{"three elements", []string{"a", "b", "c"}, []string{"b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shift(tt.a)
			if !slices.Equal(tt.want, got) {
				t.Errorf("shift() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toVersion(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s       string
		want    semverpair.Version
		wantErr bool
	}{
		{"invalid version", "blah", semverpair.Version{}, true},
		{"blank version", "", semverpair.Version{}, true},
		{"just major", "v1", semverpair.Version{Major: 1, Minor: 0, Patch: 0}, false},
		{"major and minor", "v1.2", semverpair.Version{Major: 1, Minor: 2, Patch: 0}, false},
		{"major, minor and patch", "v1.2.4", semverpair.Version{Major: 1, Minor: 2, Patch: 4}, false},
		{"major, minor, patch and dev", "v1.2.4-dev3", semverpair.Version{Major: 1, Minor: 2, Patch: 4}, false},
		{"major, minor, patch and dev+build", "v1.2.4-dev3+special", semverpair.Version{Major: 1, Minor: 2, Patch: 4}, false},
		{"major, minor, patch and build only", "v1.2.4+special", semverpair.Version{Major: 1, Minor: 2, Patch: 4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := toVersion(tt.s)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("toVersion() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("toVersion() succeeded unexpectedly")
			}
			if got.String() != tt.want.String() {
				t.Errorf("toVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodePair(t *testing.T) {
	tests := []struct {
		name    string
		decode  string
		want    semverpair.Pair
		wantErr bool
	}{
		{"blank", "", semverpair.Pair{}, true},
		{"two digits", "v1.20000.20000", semverpair.Pair{semverpair.Version{1, 0, 0}, semverpair.Version{1, 0, 0}}, false},
		{"two digits 1.0.0 & 1.0.1", "v1.20000.20001", semverpair.Pair{semverpair.Version{1, 0, 0}, semverpair.Version{1, 0, 1}}, false},
		{"three digits 1.0.0 & 1.0.1", "v1.3000000.3000001", semverpair.Pair{semverpair.Version{1, 0, 0}, semverpair.Version{1, 0, 1}}, false},
		{"two digits 1.1.0 & 1.2.1", "v1.20102.20001", semverpair.Pair{semverpair.Version{1, 1, 0}, semverpair.Version{1, 2, 1}}, false},
		{"three digits 1.1.0 & 1.2.1", "v1.3001002.3000001", semverpair.Pair{semverpair.Version{1, 1, 0}, semverpair.Version{1, 2, 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := decodePair(tt.decode)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("decodePair() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("decodePair() succeeded unexpectedly")
			}
			if got.First.String() != tt.want.First.String() || got.Second.String() != tt.want.Second.String() {
				t.Errorf("decodePair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodePair(t *testing.T) {
	tests := []struct {
		name    string
		first   string
		second  string
		want    semverpair.Pair
		wantErr bool
	}{
		{"both blank", "", "", semverpair.Pair{}, true},
		{"first blank", "", "v1.0.0", semverpair.Pair{}, true},
		{"second blank", "v1.0.0", "", semverpair.Pair{}, true},
		{"first invalid", "blah", "v1.0.0", semverpair.Pair{}, true},
		{"second invalid", "v1.0.0", "blah", semverpair.Pair{}, true},
		{"both invalid", "foo", "blah", semverpair.Pair{}, true},
		{"1.0.0 & 1.0.1", "1.0.0", "1.0.1", semverpair.Pair{semverpair.Version{1, 0, 0}, semverpair.Version{1, 0, 1}}, false},
		{"v1.0.0 & 1.0.1", "v1.0.0", "1.0.1", semverpair.Pair{semverpair.Version{1, 0, 0}, semverpair.Version{1, 0, 1}}, false},
		{"1.0.0 & v1.0.1", "1.0.0", "v1.0.1", semverpair.Pair{semverpair.Version{1, 0, 0}, semverpair.Version{1, 0, 1}}, false},
		{"v1.0.0 & v1.0.1", "v1.0.0", "v1.0.1", semverpair.Pair{semverpair.Version{1, 0, 0}, semverpair.Version{1, 0, 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := encodePair(tt.first, tt.second)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("encodePair() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("encodePair() succeeded unexpectedly")
			}
			if got.First.String() != tt.want.First.String() || got.Second.String() != tt.want.Second.String() {
				t.Errorf("encodePair() = %v, want %v", got, tt.want)
			}
		})
	}
}
