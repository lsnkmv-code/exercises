package validation

import (
	"errors"
	"os"
	"testing"

	flag "github.com/spf13/pflag"
)

func requireErrorType[T error](t *testing.T, err error) T {
	t.Helper()
	var target T
	if !errors.As(err, &target) {
		t.Fatalf("error = %T (%v), want %T", err, err, target)
	}
	return target
}

func TestValidateCommands(t *testing.T) {
	tests := []struct {
		name    string
		flags   Flags
		wantErr bool
	}{
		{name: "distance", flags: Flags{Distance: true}},
		{name: "area", flags: Flags{Area: true}},
		{name: "perimeter", flags: Flags{Perimeter: true}},
		{name: "contains", flags: Flags{Contains: true}},
		{name: "none", flags: Flags{}, wantErr: true},
		{name: "several", flags: Flags{Area: true, Perimeter: true}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommands(&tt.flags)
			if tt.wantErr {
				requireErrorType[*CommandError](t, err)
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateDistance(t *testing.T) {
	valid := Flags{Points: []string{"0", "1", "2.5", "-3"}}
	if err := validateDistance(&valid); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wrongCount := Flags{Points: []string{"0", "1"}}
	requireErrorType[*InvalidPointsCountError](t, validateDistance(&wrongCount))

	invalidNumber := Flags{Points: []string{"0", "1", "bad", "3"}}
	err := requireErrorType[*InvalidTypeError](t, validateDistance(&invalidNumber))
	if err.Got != "bad" {
		t.Fatalf("Got = %q, want bad", err.Got)
	}
}

func TestValidatePolygon(t *testing.T) {
	valid := Flags{Points: []string{"0", "0", "2", "0", "0", "2"}}
	if err := validatePolygon(&valid, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, points := range [][]string{
		{"0", "0", "1", "1"},
		{"0", "0", "1", "0", "1", "1", "0"},
	} {
		flags := Flags{Points: points}
		requireErrorType[*InvalidPointsCountError](t, validatePolygon(&flags, 0))
	}

	invalidNumber := Flags{Points: []string{"0", "0", "2", "no", "0", "2"}}
	requireErrorType[*InvalidTypeError](t, validatePolygon(&invalidNumber, 0))
}

func TestValidateCircle(t *testing.T) {
	tests := []struct {
		name      string
		flags     Flags
		errorType string
	}{
		{name: "combined", flags: Flags{Circle: []string{"1", "2", "3"}}},
		{name: "separate", flags: Flags{Center: []string{"1", "2"}, Radius: "3"}},
		{name: "missing", flags: Flags{}, errorType: "count"},
		{name: "bad combined number", flags: Flags{Circle: []string{"1", "x", "3"}}, errorType: "type"},
		{name: "bad center count", flags: Flags{Center: []string{"1"}, Radius: "3"}, errorType: "count"},
		{name: "bad radius", flags: Flags{Center: []string{"1", "2"}, Radius: "x"}, errorType: "type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCircle(&tt.flags)
			switch tt.errorType {
			case "count":
				requireErrorType[*InvalidPointsCountError](t, err)
			case "type":
				requireErrorType[*InvalidTypeError](t, err)
			default:
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateAreaOrPerimeter(t *testing.T) {
	tests := []struct {
		name    string
		flags   Flags
		wantErr bool
	}{
		{name: "polygon", flags: Flags{Polygon: true, Points: []string{"0", "0", "2", "0", "0", "2"}}},
		{name: "combined circle", flags: Flags{Circle: []string{"0", "0", "2"}}},
		{name: "separate circle", flags: Flags{Center: []string{"0", "0"}, Radius: "2"}},
		{name: "invalid polygon", flags: Flags{Polygon: true, Points: []string{"0", "0"}}, wantErr: true},
		{name: "invalid circle", flags: Flags{Circle: []string{"0", "bad", "2"}}, wantErr: true},
		{name: "missing figure", flags: Flags{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAreaOrPerimeter(&tt.flags)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateContains(t *testing.T) {
	tests := []struct {
		name      string
		flags     Flags
		errorType string
	}{
		{name: "circle", flags: Flags{Points: []string{"1", "1"}, Circle: []string{"0", "0", "2"}}},
		{name: "polygon", flags: Flags{Points: []string{"1", "1", "0", "0", "2", "0", "0", "2"}}},
		{name: "missing point", flags: Flags{Points: []string{"1"}}, errorType: "info"},
		{name: "bad point", flags: Flags{Points: []string{"x", "1"}, Circle: []string{"0", "0", "2"}}, errorType: "type"},
		{name: "invalid circle", flags: Flags{Points: []string{"1", "1"}}, errorType: "count"},
		{name: "invalid polygon", flags: Flags{Points: []string{"1", "1", "0", "0"}}, errorType: "count"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContains(&tt.flags)
			switch tt.errorType {
			case "info":
				requireErrorType[*NotEnoughInfoError](t, err)
			case "type":
				requireErrorType[*InvalidTypeError](t, err)
			case "count":
				requireErrorType[*InvalidPointsCountError](t, err)
			default:
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateFlagsDispatchesCommands(t *testing.T) {
	tests := []struct {
		name  string
		flags Flags
	}{
		{name: "distance", flags: Flags{Distance: true, Points: []string{"0", "0", "3", "4"}}},
		{name: "area", flags: Flags{Area: true, Circle: []string{"0", "0", "2"}}},
		{name: "perimeter", flags: Flags{Perimeter: true, Circle: []string{"0", "0", "2"}}},
		{name: "contains", flags: Flags{Contains: true, Points: []string{"1", "1"}, Circle: []string{"0", "0", "2"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateFlags(&tt.flags); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}

	invalid := Flags{}
	requireErrorType[*CommandError](t, validateFlags(&invalid))

	invalidCases := []Flags{
		{Distance: true, Points: []string{"0"}},
		{Area: true},
		{Perimeter: true},
		{Contains: true},
	}
	for _, flags := range invalidCases {
		if err := validateFlags(&flags); err == nil {
			t.Fatalf("validateFlags(%+v) expected an error", flags)
		}
	}
}

func TestParseFlags(t *testing.T) {
	originalArgs := os.Args
	originalCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = originalArgs
		flag.CommandLine = originalCommandLine
	})

	flag.CommandLine = flag.NewFlagSet("geometry-test", flag.ContinueOnError)
	os.Args = []string{"geometry", "--distance", "--point", "0,0,3,4"}

	flags, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}
	if !flags.Distance || len(flags.Points) != 4 {
		t.Fatalf("unexpected flags: %+v", flags)
	}
}

func TestParseFlagsValidationError(t *testing.T) {
	originalArgs := os.Args
	originalCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = originalArgs
		flag.CommandLine = originalCommandLine
	})

	flag.CommandLine = flag.NewFlagSet("geometry-test-error", flag.ContinueOnError)
	os.Args = []string{"geometry"}

	_, err := ParseFlags()
	requireErrorType[*CommandError](t, err)
}
