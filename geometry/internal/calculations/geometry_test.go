package calculations

import (
	"geometry/cmd/app/internal/validation"
	"math"
	"strings"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestPointDistanceTo(t *testing.T) {
	point := Point{X: 1, Y: 2}
	if got := point.DistanceTo(Point{X: 4, Y: 6}); got != 5 {
		t.Fatalf("DistanceTo() = %v, want 5", got)
	}
}

func TestPolygonAreaAndPerimeter(t *testing.T) {
	polygon := Polygon{Points: []Point{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 3}}}

	if got := polygon.Area(); got != 6 {
		t.Fatalf("Area() = %v, want 6", got)
	}
	if got := polygon.Perimeter(); got != 12 {
		t.Fatalf("Perimeter() = %v, want 12", got)
	}

	reversed := Polygon{Points: []Point{{X: 4, Y: 3}, {X: 4, Y: 0}, {X: 0, Y: 0}}}
	if got := reversed.Area(); got != 6 {
		t.Fatalf("Area() for clockwise polygon = %v, want 6", got)
	}
}

func TestPolygonContains(t *testing.T) {
	polygon := Polygon{Points: []Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}}

	tests := []struct {
		name  string
		point Point
		want  bool
	}{
		{name: "inside", point: Point{X: 5, Y: 5}, want: true},
		{name: "outside", point: Point{X: 15, Y: 5}, want: false},
		{name: "on horizontal edge", point: Point{X: 5, Y: 0}, want: true},
		{name: "on vertical edge", point: Point{X: 10, Y: 5}, want: true},
		{name: "on vertex", point: Point{X: 0, Y: 0}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := polygon.Contains(tt.point); got != tt.want {
				t.Fatalf("Contains(%v) = %v, want %v", tt.point, got, tt.want)
			}
		})
	}
}

func TestPolygonContainsConcavePolygon(t *testing.T) {
	polygon := Polygon{Points: []Point{
		{X: 0, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 4},
		{X: 2, Y: 2},
		{X: 0, Y: 4},
	}}

	if !polygon.Contains(Point{X: 1, Y: 1}) {
		t.Fatal("expected point in polygon to be contained")
	}
	if polygon.Contains(Point{X: 2, Y: 3}) {
		t.Fatal("expected point in concave cutout to be outside")
	}
}

func TestPolygonContainsInvalidPolygon(t *testing.T) {
	polygon := Polygon{Points: []Point{{X: 0, Y: 0}, {X: 1, Y: 1}}}
	if polygon.Contains(Point{}) {
		t.Fatal("polygon with fewer than three vertices must not contain a point")
	}
}

func TestCircleCalculations(t *testing.T) {
	circle := Circle{Center: Point{X: 1, Y: 1}, Radius: 2}

	if got := circle.Area(); !almostEqual(got, 4*math.Pi) {
		t.Fatalf("Area() = %v, want %v", got, 4*math.Pi)
	}
	if got := circle.Perimeter(); !almostEqual(got, 4*math.Pi) {
		t.Fatalf("Perimeter() = %v, want %v", got, 4*math.Pi)
	}
	if !circle.Contains(Point{X: 3, Y: 1}) {
		t.Fatal("point on circle boundary must be contained")
	}
	if circle.Contains(Point{X: 3.1, Y: 1}) {
		t.Fatal("point outside circle must not be contained")
	}
}

func TestCreateFigures(t *testing.T) {
	tests := []struct {
		name  string
		flags validation.Flags
		check func(t *testing.T, figures Figures)
	}{
		{
			name:  "distance",
			flags: validation.Flags{Distance: true, Points: []string{"0", "1", "2", "3"}},
			check: func(t *testing.T, figures Figures) {
				if len(figures.Points) != 2 || figures.Points[1] != (Point{X: 2, Y: 3}) {
					t.Fatalf("unexpected points: %+v", figures.Points)
				}
			},
		},
		{
			name:  "polygon area",
			flags: validation.Flags{Area: true, Polygon: true, Points: []string{"0", "0", "2", "0", "0", "2"}},
			check: func(t *testing.T, figures Figures) {
				if len(figures.Polygons.Points) != 3 {
					t.Fatalf("unexpected polygon: %+v", figures.Polygons)
				}
			},
		},
		{
			name:  "circle perimeter from combined flag",
			flags: validation.Flags{Perimeter: true, Circle: []string{"1", "2", "3"}},
			check: func(t *testing.T, figures Figures) {
				if figures.Circles != (Circle{Center: Point{X: 1, Y: 2}, Radius: 3}) {
					t.Fatalf("unexpected circle: %+v", figures.Circles)
				}
			},
		},
		{
			name: "circle area from separate flags",
			flags: validation.Flags{
				Area: true, Center: []string{"4", "5"}, Radius: "6",
			},
			check: func(t *testing.T, figures Figures) {
				if figures.Circles != (Circle{Center: Point{X: 4, Y: 5}, Radius: 6}) {
					t.Fatalf("unexpected circle: %+v", figures.Circles)
				}
			},
		},
		{
			name: "polygon contains",
			flags: validation.Flags{
				Contains: true, Polygon: true,
				Points: []string{"1", "1", "0", "0", "2", "0", "0", "2"},
			},
			check: func(t *testing.T, figures Figures) {
				if figures.Points[0] != (Point{X: 1, Y: 1}) || len(figures.Polygons.Points) != 3 {
					t.Fatalf("unexpected figures: %+v", figures)
				}
			},
		},
		{
			name: "circle contains",
			flags: validation.Flags{
				Contains: true, Points: []string{"1", "1"}, Circle: []string{"0", "0", "2"},
			},
			check: func(t *testing.T, figures Figures) {
				if figures.Points[0] != (Point{X: 1, Y: 1}) || figures.Circles.Radius != 2 {
					t.Fatalf("unexpected figures: %+v", figures)
				}
			},
		},
		{
			name:  "no operation",
			flags: validation.Flags{},
			check: func(t *testing.T, figures Figures) {
				if len(figures.Points) != 0 || len(figures.Polygons.Points) != 0 || figures.Circles.Radius != 0 {
					t.Fatalf("expected empty figures, got %+v", figures)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, CreateFigures(&tt.flags))
		})
	}
}

func TestFiguresCalculate(t *testing.T) {
	triangle := Polygon{Points: []Point{{X: 0, Y: 0}, {X: 3, Y: 0}, {X: 0, Y: 4}}}
	tests := []struct {
		name    string
		figures Figures
		flags   validation.Flags
		want    string
	}{
		{name: "distance", figures: Figures{Points: []Point{{}, {X: 3, Y: 4}}}, flags: validation.Flags{Distance: true}, want: "The distance is: 5.00"},
		{name: "polygon area", figures: Figures{Polygons: triangle}, flags: validation.Flags{Area: true, Polygon: true}, want: "The area of given polygon is: 6.00"},
		{name: "circle area", figures: Figures{Circles: Circle{Radius: 2}}, flags: validation.Flags{Area: true}, want: "The area of given circle is: 12.57"},
		{name: "polygon perimeter", figures: Figures{Polygons: triangle}, flags: validation.Flags{Perimeter: true, Polygon: true}, want: "The perimeter of given polygon is: 12.00"},
		{name: "circle perimeter", figures: Figures{Circles: Circle{Radius: 2}}, flags: validation.Flags{Perimeter: true}, want: "The perimeter of given circle is: 12.57"},
		{name: "polygon contains", figures: Figures{Points: []Point{{X: 1, Y: 1}}, Polygons: triangle}, flags: validation.Flags{Contains: true, Polygon: true}, want: "Does the polygon contain the point: true"},
		{name: "circle contains", figures: Figures{Points: []Point{{X: 1, Y: 1}}, Circles: Circle{Radius: 2}}, flags: validation.Flags{Contains: true}, want: "Does the circle contain the point: true"},
		{name: "no operation", figures: Figures{}, flags: validation.Flags{}, want: "something went wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.figures.Calculate(&tt.flags); got != tt.want {
				t.Fatalf("Calculate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCalculateResultContainsTwoDecimals(t *testing.T) {
	figures := Figures{Circles: Circle{Radius: 1}}
	flags := validation.Flags{Area: true}
	if got := figures.Calculate(&flags); !strings.HasSuffix(got, "3.14") {
		t.Fatalf("expected two decimal places, got %q", got)
	}
}
