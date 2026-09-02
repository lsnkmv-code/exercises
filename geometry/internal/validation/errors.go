package validation

import "fmt"

type CommandError struct {
	NumberOfCommands int
	ExpectedNumber   int
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("number of given commands: %d; number of expected commands: %d", e.NumberOfCommands, e.ExpectedNumber)
}

type InvalidPointsCountError struct {
	TypeOfPoint          string
	NumberOfNeededPoints string
	Coordinates          string
}

func (e *InvalidPointsCountError) Error() string {
	return fmt.Sprintf("invalid %s count, expected: %s points with %s coordinates", e.TypeOfPoint, e.NumberOfNeededPoints, e.Coordinates)
}

type InvalidTypeError struct {
	Parameter string
	Expected  string
	Got       string
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type of %s, expected: %s, got parameter with value: %s", e.Parameter, e.Expected, e.Got)
}

type NotEnoughInfoError struct{
	Type string
}

func (e *NotEnoughInfoError) Error() string{
	return fmt.Sprintf("not enough parameters given, missing %s", e.Type)
}