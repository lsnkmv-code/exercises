package validation

import (
	"strconv"

	flag "github.com/spf13/pflag"
)

type Flags struct {
	Distance  bool
	Area      bool
	Perimeter bool
	Contains  bool
	Points    []string
	Circle    []string
	Center    []string
	Radius    string
	Polygon   bool
}

func ParseFlags() (Flags, error) {
	var flags Flags

	flag.BoolVar(&flags.Distance, "distance", false, "One of the cases when user wants to calculate the distance from one point to another")
	flag.BoolVar(&flags.Area, "area", false, "One of the cases when user wants to calculate the area of a figure")
	flag.BoolVar(&flags.Perimeter, "perimeter", false, "One of the cases when user wants to calculate the perimeter of a figure")
	flag.BoolVar(&flags.Contains, "contains", false, "One of the cases when user wants to know if certain point is contained in a certain figure")
	flag.StringSliceVar(&flags.Points, "point", nil, "point coordinates X, Y")
	flag.StringSliceVar(&flags.Circle, "circle", nil, "circle center X, Y coordinates and radius R (X, Y, R)")
	flag.StringSliceVar(&flags.Center, "center", nil, "center coordinates X, Y")
	flag.StringVar(&flags.Radius, "radius", "", "radius")
	flag.BoolVar(&flags.Polygon, "polygon", false, "Is figure a polygon")

	flag.Parse()
	if err := validateFlags(&flags); err != nil {
		return flags, err
	}
	return flags, nil
}

func validateFlags(flags *Flags) error {
	if err := validateCommands(flags); err != nil {
		return err
	}
	if flags.Distance {
		if err := validateDistance(flags); err != nil {
			return err
		}
	}
	if flags.Area || flags.Perimeter {
		if err := validateAreaOrPerimeter(flags); err != nil {
			return err
		}
	}
	if flags.Contains {
		if err := validateContains(flags); err != nil {
			return err
		}
	}
	return nil
}

func validateCommands(flags *Flags) error {
	countCommands := 0

	if flags.Distance {
		countCommands += 1
	}
	if flags.Area {
		countCommands += 1
	}
	if flags.Contains {
		countCommands += 1
	}
	if flags.Perimeter {
		countCommands += 1
	}

	if countCommands == 0 || countCommands > 1 {
		return &CommandError{
			NumberOfCommands: countCommands,
			ExpectedNumber:   1,
		}
	}
	return nil
}

func validateDistance(flags *Flags) error {
	if len(flags.Points) != 4 {
		return &InvalidPointsCountError{
			TypeOfPoint:          "points for distance",
			NumberOfNeededPoints: "2",
			Coordinates:          "X,Y",
		}
	}
	for _, point := range flags.Points {
		if _, err := strconv.ParseFloat(point, 64); err != nil {
			return &InvalidTypeError{
				Parameter: "point for distance",
				Expected:  "float64",
				Got:       point,
			}
		}
	}

	return nil
}

func validateAreaOrPerimeter(flags *Flags) error {
	if flags.Polygon {
		if err := validatePolygon(flags, 0); err != nil {
			return err
		}
	} else if len(flags.Circle) != 0 || flags.Center != nil && flags.Radius != "" {
		if err := validateCircle(flags); err != nil {
			return err
		}
	} else {
		return &NotEnoughInfoError{
			Type: "circle or polygon",
		}
	}
	return nil
}

func validatePolygon(flags *Flags, point int) error {
	if len(flags.Points) < (6+point) || (len(flags.Points)-point)%2 != 0 {
		return &InvalidPointsCountError{
			TypeOfPoint:          "points for polygon",
			NumberOfNeededPoints: "at least 3",
			Coordinates:          "X,Y",
		}
	}
	for _, point := range flags.Points {
		if _, err := strconv.ParseFloat(point, 64); err != nil {
			return &InvalidTypeError{
				Parameter: "point for polygon",
				Expected:  "float64",
				Got:       point,
			}
		}
	}
	return nil
}

func validateCircle(flags *Flags) error {
	if len(flags.Circle) != 3 && flags.Center == nil && flags.Radius == "" {
		return &InvalidPointsCountError{
			TypeOfPoint:          "points for circle",
			NumberOfNeededPoints: "3",
			Coordinates:          "X,Y,R",
		}
	}
	if len(flags.Circle) == 3 {
		for _, point := range flags.Circle {
			if _, err := strconv.ParseFloat(point, 64); err != nil {
				return &InvalidTypeError{
					Parameter: "circle coordinates",
					Expected:  "float64",
					Got:       point,
				}
			}
		}
	}

	if len(flags.Circle) == 0 && len(flags.Center) != 2 {
		return &InvalidPointsCountError{
			TypeOfPoint:          "points for circle center",
			NumberOfNeededPoints: "1",
			Coordinates:          "X,Y",
		}
	}
	if len(flags.Circle) == 0 && len(flags.Center) == 2 {
		if r, err := strconv.ParseFloat(flags.Radius, 64); err != nil || r < 0 {
			return &InvalidTypeError{
				Parameter: "circle radius (positive number)",
				Expected:  "float64",
				Got:       flags.Radius,
			}
		}
	}

	return nil
}

func validateContains(flags *Flags) error {
	if len(flags.Points) < 2 {
		return &NotEnoughInfoError{
			Type: "points for contains command",
		}
	} else {
		for _, point := range flags.Points[0:2] {
			if _, err := strconv.ParseFloat(point, 64); err != nil {
				return &InvalidTypeError{
					Parameter: "point (contains)",
					Expected:  "float64",
					Got:       point,
				}
			}
		}
	}

	if len(flags.Points) == 2 {
		if err := validateCircle(flags); err != nil {
			return err
		}
	} else {
		if err := validatePolygon(flags, 2); err != nil {
			return err
		}
	}

	return nil
}
