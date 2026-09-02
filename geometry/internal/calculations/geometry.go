package calculations

import (
	"fmt"
	"geometry/cmd/app/internal/validation"
	"math"
	"strconv"
)

type Figures struct {
	Points   []Point
	Circles  Circle
	Polygons Polygon
}

type Point struct {
	X float64
	Y float64
}

func (p *Point) DistanceTo(p2 Point) float64 {
	return math.Hypot(p.X-p2.X, p.Y-p2.Y)
}

type Polygon struct {
	Points []Point
}

func (polygon *Polygon) Area() float64 {
	ans := 0.0
	for i := range polygon.Points {
		if i == len(polygon.Points)-1 {
			ans += polygon.Points[i].X*polygon.Points[0].Y - polygon.Points[0].X*polygon.Points[i].Y
		} else {
			ans += polygon.Points[i].X*polygon.Points[i+1].Y - polygon.Points[i+1].X*polygon.Points[i].Y
		}
	}

	return math.Abs(ans / 2)
}
func (polygon *Polygon) Perimeter() float64 {
	ans := 0.0
	for i := range polygon.Points {
		if i == len(polygon.Points)-1 {
			ans += polygon.Points[i].DistanceTo(polygon.Points[0])
		} else {
			ans += polygon.Points[i].DistanceTo(polygon.Points[i+1])
		}
	}

	return ans
}

func (polygon *Polygon) Contains(point Point) bool {
	if len(polygon.Points) < 3 {
		return false
	}

	inside := false
	previous := polygon.Points[len(polygon.Points)-1]

	for _, current := range polygon.Points {
		if pointOnSegment(point, previous, current) {
			return true
		}
		if (previous.Y > point.Y) != (current.Y > point.Y) {
			intersectionX := previous.X +
				(point.Y-previous.Y)*(current.X-previous.X)/(current.Y-previous.Y)
			if point.X < intersectionX {
				inside = !inside
			}
		}

		previous = current
	}

	return inside
}

func pointOnSegment(point, start, end Point) bool {
	const epsilon = 1e-9

	crossProduct := (point.X-start.X)*(end.Y-start.Y) -
		(point.Y-start.Y)*(end.X-start.X)
	if math.Abs(crossProduct) > epsilon {
		return false
	}

	dotProduct := (point.X-start.X)*(point.X-end.X) +
		(point.Y-start.Y)*(point.Y-end.Y)
	return dotProduct <= epsilon
}

type Circle struct {
	Center Point
	Radius float64
}

func (circle *Circle) Area() float64 {
	return math.Pi * circle.Radius * circle.Radius
}

func (circle *Circle) Perimeter() float64 {
	return 2 * math.Pi * circle.Radius
}

func (circle *Circle) Contains(point Point) bool {
	dx := point.X - circle.Center.X
	dy := point.Y - circle.Center.Y
	return (dx*dx + dy*dy) <= (circle.Radius * circle.Radius)
}

func CreateFigures(flags *validation.Flags) Figures {
	var parsedPoints []float64
	for _, point := range flags.Points {
		val, _ := strconv.ParseFloat(point, 64)
		parsedPoints = append(parsedPoints, val)
	}

	if flags.Distance {
		var points []Point
		for i := 0; i < 4; i += 2 {
			points = append(points, Point{X: parsedPoints[i], Y: parsedPoints[i+1]})
		}

		return Figures{Points: points}
	}

	if flags.Area || flags.Perimeter {
		if flags.Polygon {
			polygon := createPolygon(0, parsedPoints)
			return Figures{Polygons: polygon}
		} else {
			circle := createCircle(flags)
			return Figures{Circles: circle}
		}
	}

	if flags.Contains {
		points := []Point{{X: parsedPoints[0], Y: parsedPoints[1]}}
		if flags.Polygon {
			polygon := createPolygon(2, parsedPoints)
			return Figures{Points: points, Polygons: polygon}
		} else {
			circle := createCircle(flags)
			return Figures{Points: points, Circles: circle}
		}
	}

	return Figures{}
}

func createPolygon(start int, parsedPoints []float64) Polygon {
	var points []Point
	for i := start; i < len(parsedPoints); i += 2 {
		points = append(points, Point{X: parsedPoints[i], Y: parsedPoints[i+1]})
	}
	return Polygon{Points: points}
}

func createCircle(flags *validation.Flags) Circle {
	if flags.Circle != nil {
		xCenter, _ := strconv.ParseFloat(flags.Circle[0], 64)
		yCenter, _ := strconv.ParseFloat(flags.Circle[1], 64)
		radius, _ := strconv.ParseFloat(flags.Circle[2], 64)
		circle := Circle{Center: Point{X: xCenter, Y: yCenter}, Radius: radius}

		return circle
	} else {
		xCenter, _ := strconv.ParseFloat(flags.Center[0], 64)
		yCenter, _ := strconv.ParseFloat(flags.Center[1], 64)
		radius, _ := strconv.ParseFloat(flags.Radius, 64)
		circle := Circle{Center: Point{X: xCenter, Y: yCenter}, Radius: radius}

		return circle
	}
}

func (f *Figures) Calculate(flags *validation.Flags) string {
	if flags.Distance {
		distance := f.Points[0].DistanceTo(f.Points[1])
		return fmt.Sprintf("The distance is: %.2f", distance)
	}

	if flags.Area {
		if flags.Polygon {
			area := f.Polygons.Area()
			return fmt.Sprintf("The area of given polygon is: %.2f", area)
		} else {
			area := f.Circles.Area()
			return fmt.Sprintf("The area of given circle is: %.2f", area)
		}
	}
	if flags.Perimeter {
		if flags.Polygon {
			area := f.Polygons.Perimeter()
			return fmt.Sprintf("The perimeter of given polygon is: %.2f", area)
		} else {
			area := f.Circles.Perimeter()
			return fmt.Sprintf("The perimeter of given circle is: %.2f", area)
		}
	}
	if flags.Contains {
		if flags.Polygon {
			contains := f.Polygons.Contains(f.Points[0])
			return fmt.Sprintf("Does the polygon contain the point: %v", contains)
		} else {
			contains := f.Circles.Contains(f.Points[0])
			return fmt.Sprintf("Does the circle contain the point: %v", contains)
		}
	}
	return "something went wrong"
}
