package geo

import (
	"math"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS = 6371
)

// Geometry is a geometric shape that can be used to
// verify whether Geo point is contained in there or not.
type Geometry interface {
	Contains(latlng LatLng) bool

	CircleBound() (LatLng, float64)
}

type LatLng struct {
	Lat float64
	Lng float64
}

// GreatCircleDistance: Calculates the Haversine distance between two points in kilometers.
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *LatLng) GreatCircleDistance(p2 *LatLng) float64 {
	dLat := (p2.Lat - p.Lat) * (math.Pi / 180.0)
	dLon := (p2.Lng - p.Lng) * (math.Pi / 180.0)

	lat1 := p.Lat * (math.Pi / 180.0)
	lat2 := p2.Lat * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}

type Point struct {
	Type        string
	Coordinates []float64
}

func (p *Point) Contains(latlng LatLng) bool {
	return latlng.Lng == p.Coordinates[0] && latlng.Lat == p.Coordinates[1]
}

func (p *Point) CircleBound() (LatLng, float64) {
	return LatLng{Lng: p.Coordinates[0], Lat: p.Coordinates[1]}, 0
}

// Polygon is representing a polygon
// line structure.
type Polygon struct {
	Type        string
	coordinates [][][]float64

	loop *s2.Loop
}

// NewPolygon creates a new polygon for the provided
// coordinates.
func NewPolygon(coordinates [][][]float64) *Polygon {
	points := make([]s2.Point, len(coordinates[0]))

	for i, coordinate := range coordinates[0] {
		points[len(coordinates[0])-1-i] = s2.PointFromLatLng(s2.LatLngFromDegrees(coordinate[1], coordinate[0]))
	}

	loop := s2.LoopFromPoints(points)
	return &Polygon{coordinates: coordinates, loop: loop}
}

// Contains checks whether the LatLng is contained in the
// polygon. It returns true if the LatLng is contained in
// the polygon and false otherwise.
func (p *Polygon) Contains(latlng LatLng) bool {
	return p.loop.ContainsPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(latlng.Lat, latlng.Lng)))
}

func (p *Polygon) CircleBound() (LatLng, float64) {
	points := make([]s2.Point, len(p.coordinates[0]))

	for i, coordinate := range p.coordinates[0] {
		points[len(p.coordinates[0])-1-i] = s2.PointFromLatLng(s2.LatLngFromDegrees(coordinate[1], coordinate[0]))
	}

	cap := s2.LoopFromPoints(points).CapBound()
	latlng := s2.LatLngFromPoint(cap.Center())

	return LatLng{
		Lng: latlng.Lng.Degrees(),
		Lat: latlng.Lat.Degrees(),
	}, float64(cap.Radius() * 6371000)
}

type Circle struct {
	Type        string
	Radius      float64
	Coordinates []float64
}

func (c *Circle) Contains(latlng LatLng) bool {
	earthRadius := float64(6371000)
	radiusAngle := s1.Angle(c.Radius / earthRadius)

	a := s2.PointFromLatLng(s2.LatLngFromDegrees(c.Coordinates[1], c.Coordinates[0]))
	b := s2.PointFromLatLng(s2.LatLngFromDegrees(latlng.Lat, latlng.Lng))

	rd := s1.ChordAngleFromAngle(radiusAngle)
	cp := s2.CapFromCenterChordAngle(a, rd)

	return cp.ContainsPoint(b)
}

func (c *Circle) CircleBound() (LatLng, float64) {
	return LatLng{Lng: c.Coordinates[0], Lat: c.Coordinates[1]}, c.Radius
}

type Rectangle struct {
	Type        string
	Coordinates [][]float64
}

func (r *Rectangle) Contains(latlng LatLng) bool {
	rect := s2.RectFromLatLng(s2.LatLngFromDegrees(r.Coordinates[0][1], r.Coordinates[0][0]))

	for i := 1; i < len(r.Coordinates[0]); i++ {
		rect = rect.AddPoint(s2.LatLngFromDegrees(r.Coordinates[i][1], r.Coordinates[i][0]))
	}

	return rect.ContainsLatLng(s2.LatLngFromDegrees(latlng.Lat, latlng.Lng))
}

func (r *Rectangle) CircleBound() (LatLng, float64) {
	rect := s2.RectFromLatLng(s2.LatLngFromDegrees(r.Coordinates[0][1], r.Coordinates[0][0]))

	for i := 1; i < len(r.Coordinates[0]); i++ {
		rect = rect.AddPoint(s2.LatLngFromDegrees(r.Coordinates[i][1], r.Coordinates[i][0]))
	}

	cap := rect.CapBound()
	latlng := s2.LatLngFromPoint(cap.Center())

	return LatLng{
		Lat: latlng.Lat.Degrees(),
		Lng: latlng.Lng.Degrees(),
	}, float64(cap.Radius() * 6371000)
}
