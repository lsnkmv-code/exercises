package validation

import "testing"

func TestValidationErrorMessages(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "command",
			err:  &CommandError{NumberOfCommands: 2, ExpectedNumber: 1},
			want: "number of given commands: 2; number of expected commands: 1",
		},
		{
			name: "points count",
			err: &InvalidPointsCountError{
				TypeOfPoint: "vertices", NumberOfNeededPoints: "3", Coordinates: "X,Y",
			},
			want: "invalid vertices count, expected: 3 points with X,Y coordinates",
		},
		{
			name: "type",
			err:  &InvalidTypeError{Parameter: "radius", Expected: "float64", Got: "abc"},
			want: "invalid type of radius, expected: float64, got parameter with value: abc",
		},
		{
			name: "not enough info",
			err:  &NotEnoughInfoError{Type: "circle or polygon"},
			want: "not enough parameters given, missing circle or polygon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
