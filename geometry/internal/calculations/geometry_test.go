package calculations

import (
	"geometry/cmd/app/internal/validation"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPointDistanceTo(t *testing.T) {
	point := Point{X: 1, Y: 2}
	assert.Equal(t, 5.0, point.DistanceTo(Point{X: 4, Y: 6}))
}

func TestPolygonAreaAndPerimeter(t *testing.T) {
	polygon := Polygon{Points: []Point{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 3}}}

	assert.Equal(t, 6.0, polygon.Area())
	assert.Equal(t, 12.0, polygon.Perimeter())

	reversed := Polygon{Points: []Point{{X: 4, Y: 3}, {X: 4, Y: 0}, {X: 0, Y: 0}}}
	assert.Equal(t, 6.0, reversed.Area())
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
			assert.Equal(t, tt.want, polygon.Contains(tt.point))
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

	assert.True(t, polygon.Contains(Point{X: 1, Y: 1}))
	assert.False(t, polygon.Contains(Point{X: 2, Y: 3}))
}

func TestPolygonContainsInvalidPolygon(t *testing.T) {
	polygon := Polygon{Points: []Point{{X: 0, Y: 0}, {X: 1, Y: 1}}}
	assert.False(t, polygon.Contains(Point{}))
}

func TestCircleCalculations(t *testing.T) {
	circle := Circle{Center: Point{X: 1, Y: 1}, Radius: 2}

	assert.InDelta(t, 4*math.Pi, circle.Area(), 1e-9)
	assert.InDelta(t, 4*math.Pi, circle.Perimeter(), 1e-9)
	assert.True(t, circle.Contains(Point{X: 3, Y: 1}))
	assert.False(t, circle.Contains(Point{X: 3.1, Y: 1}))
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
				require.Len(t, figures.Points, 2)
				assert.Equal(t, Point{X: 2, Y: 3}, figures.Points[1])
			},
		},
		{
			name:  "polygon area",
			flags: validation.Flags{Area: true, Polygon: true, Points: []string{"0", "0", "2", "0", "0", "2"}},
			check: func(t *testing.T, figures Figures) {
				assert.Len(t, figures.Polygons.Points, 3)
			},
		},
		{
			name:  "circle perimeter from combined flag",
			flags: validation.Flags{Perimeter: true, Circle: []string{"1", "2", "3"}},
			check: func(t *testing.T, figures Figures) {
				assert.Equal(t, Circle{Center: Point{X: 1, Y: 2}, Radius: 3}, figures.Circles)
			},
		},
		{
			name: "circle area from separate flags",
			flags: validation.Flags{
				Area: true, Center: []string{"4", "5"}, Radius: "6",
			},
			check: func(t *testing.T, figures Figures) {
				assert.Equal(t, Circle{Center: Point{X: 4, Y: 5}, Radius: 6}, figures.Circles)
			},
		},
		{
			name: "polygon contains",
			flags: validation.Flags{
				Contains: true, Polygon: true,
				Points: []string{"1", "1", "0", "0", "2", "0", "0", "2"},
			},
			check: func(t *testing.T, figures Figures) {
				require.NotEmpty(t, figures.Points)
				assert.Equal(t, Point{X: 1, Y: 1}, figures.Points[0])
				assert.Len(t, figures.Polygons.Points, 3)
			},
		},
		{
			name: "circle contains",
			flags: validation.Flags{
				Contains: true, Points: []string{"1", "1"}, Circle: []string{"0", "0", "2"},
			},
			check: func(t *testing.T, figures Figures) {
				require.NotEmpty(t, figures.Points)
				assert.Equal(t, Point{X: 1, Y: 1}, figures.Points[0])
				assert.Equal(t, 2.0, figures.Circles.Radius)
			},
		},
		{
			name:  "no operation",
			flags: validation.Flags{},
			check: func(t *testing.T, figures Figures) {
				assert.Empty(t, figures.Points)
				assert.Empty(t, figures.Polygons.Points)
				assert.Zero(t, figures.Circles.Radius)
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
			assert.Equal(t, tt.want, tt.figures.Calculate(&tt.flags))
		})
	}
}

func TestCalculateResultContainsTwoDecimals(t *testing.T) {
	figures := Figures{Circles: Circle{Radius: 1}}
	flags := validation.Flags{Area: true}
	assert.Equal(t, "The area of given circle is: 3.14", figures.Calculate(&flags))
}
