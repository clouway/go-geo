package geo

import (
	"io"
	"os"
	"reflect"
	"testing"

	geojson "github.com/paulmach/go.geojson"
)

func TestPointContains(t *testing.T) {
	point := Point{
		Coordinates: []float64{25.608530044555668, 43.07380969664719},
	}

	cases := []struct {
		in   LatLng
		want bool
	}{
		{
			LatLng{Lat: 43.07380969664719, Lng: 25.608530044555668},
			true,
		},
		{
			LatLng{Lat: 43.11111111111111, Lng: 25.608530044555668},
			false,
		},
		{
			LatLng{Lat: 43.07380969664719, Lng: 25.222222222222222},
			false,
		},
	}

	for _, c := range cases {
		got := point.Contains(c.in)

		if got != c.want {
			t.Errorf("expected %t", c.want)
			t.Errorf("     got %t", got)
		}
	}
}

func TestPointCircleBound(t *testing.T) {
	point := Point{
		Coordinates: []float64{25.608530044555668, 43.07380969664719},
	}

	center, radius := point.CircleBound()

	wantRadius := 0.0

	wantCenter := LatLng{
		Lat: 43.07380969664719,
		Lng: 25.608530044555668,
	}

	if radius != wantRadius {
		t.Errorf("expected radius %v", wantRadius)
		t.Errorf("     got radius %v", radius)
	}

	if !reflect.DeepEqual(center, wantCenter) {
		t.Errorf("expected center %v", wantCenter)
		t.Errorf("     got center %v", center)
	}
}

func TestPolygonContains(t *testing.T) {
	polygon := NewPolygon([][][]float64{{
		{25.7244873046875, 43.11110313559475},
		{25.726847648620605, 43.11417334786724},
		{25.73268413543701, 43.110163243903585},
		{25.728735923767093, 43.10712416198819},
		{25.724401473999023, 43.10865938717618},
		{25.7244873046875, 43.11110313559475},
	}},
	)

	cases := []struct {
		in   LatLng
		want bool
	}{
		{
			LatLng{Lat: 43.1089613, Lng: 25.7267396},
			true,
		},
		{
			LatLng{Lat: 43.0765023, Lng: 25.6312193},
			false,
		},
	}

	for _, c := range cases {
		got := polygon.Contains(c.in)

		if got != c.want {
			t.Errorf("expected %t", c.want)
			t.Errorf("     got %t", got)
		}
	}
}

func TestCountryPolygonContains(t *testing.T) {
	cp := readCountryGeoJson()

	cases := []struct {
		in   LatLng
		want bool
	}{
		{
			// In the country
			LatLng{Lat: 43.1089613, Lng: 25.7267396},
			true,
		},
		{
			// In the country (Close to the borders)
			LatLng{Lat: 43.830041, Lng: 25.937221},
			true,
		},
		{
			// Romania (should be out of borders)
			LatLng{Lat: 43.855333, Lng: 25.915012},
			false,
		},
		{
			// Turkey (should be out of borders)
			LatLng{Lat: 41.677080, Lng: 26.545086},
			false,
		},
	}

	polygon := NewPolygon(cp.Features[0].Geometry.Polygon)
	for _, c := range cases {
		got := polygon.Contains(c.in)

		if got != c.want {
			t.Errorf("expected %t", c.want)
			t.Errorf("     got %t", got)
		}
	}
}

func TestPolygonCircleBound(t *testing.T) {
	polygon := NewPolygon([][][]float64{{
		{25.7244873046875, 43.11110313559475},
		{25.726847648620605, 43.11417334786724},
		{25.73268413543701, 43.110163243903585},
		{25.728735923767093, 43.10712416198819},
		{25.724401473999023, 43.10865938717618},
		{25.7244873046875, 43.11110313559475},
	}},
	)

	center, radius := polygon.CircleBound()

	wantRadius := 516.3532563112891

	wantCenter := LatLng{
		Lat: 43.11064875492771,
		Lng: 25.728542804718014,
	}

	if radius != wantRadius {
		t.Errorf("expected radius %v", wantRadius)
		t.Errorf("     got radius %v", radius)
	}

	if !reflect.DeepEqual(center, wantCenter) {
		t.Errorf("expected center %v", wantCenter)
		t.Errorf("     got center %v", center)
	}
}

func TestCircleContains(t *testing.T) {
	circle := Circle{
		Radius:      297.82929433609627,
		Coordinates: []float64{25.608530044555668, 43.07380969664719},
	}

	cases := []struct {
		in   LatLng
		want bool
	}{
		{
			LatLng{Lat: 43.07415, Lng: 25.61671},
			false,
		},
		{
			LatLng{Lat: 43.07409, Lng: 25.60987},
			true,
		},
	}

	for _, c := range cases {
		got := circle.Contains(c.in)

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("expected %t", c.want)
			t.Errorf("     got %t", got)
		}
	}
}

func TestCircleBound(t *testing.T) {
	circle := Circle{
		Radius:      297.82929433609627,
		Coordinates: []float64{25.608530044555668, 43.07380969664719},
	}

	center, radius := circle.CircleBound()

	wantRadius := 297.82929433609627

	wantCenter := LatLng{
		Lat: 43.07380969664719,
		Lng: 25.608530044555668,
	}

	if radius != wantRadius {
		t.Errorf("expected radius %v", wantRadius)
		t.Errorf("     got radius %v", radius)
	}

	if !reflect.DeepEqual(center, wantCenter) {
		t.Errorf("expected center %v", wantCenter)
		t.Errorf("     got center %v", center)
	}
}

func TestRectangleContains(t *testing.T) {
	rectangle := Rectangle{Coordinates: [][]float64{
		{25.288888, 42.244444},
		{25.322222, 42.288888},
	}}

	cases := []struct {
		in   LatLng
		want bool
	}{
		{
			LatLng{Lat: 42.266667, Lng: 25.305549},
			true,
		},
		{
			LatLng{Lat: 42.33029, Lng: 25.22495},
			false,
		},
	}

	for _, c := range cases {
		got := rectangle.Contains(c.in)

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("expected %t", c.want)
			t.Errorf("     got %t", got)
		}
	}
}

func TestRectangleCircleBound(t *testing.T) {
	rectangle := Rectangle{Coordinates: [][]float64{
		{25.288888, 42.244444},
		{25.322222, 42.288888},
	}}

	center, radius := rectangle.CircleBound()

	wantRadius := 2826.1834095674053

	wantCenter := LatLng{
		Lat: 42.266666,
		Lng: 25.305555000000002,
	}

	if radius != wantRadius {
		t.Errorf("expected radius %v", wantRadius)
		t.Errorf("     got radius %v", radius)
	}

	if !reflect.DeepEqual(center, wantCenter) {
		t.Errorf("expected center %v", wantCenter)
		t.Errorf("     got center %v", center)
	}
}

func BenchmarkPolygonContains(b *testing.B) {
	polygon := NewPolygon(
		[][][]float64{{
			{25.7244873046875, 43.11110313559475},
			{25.726847648620605, 43.11417334786724},
			{25.73268413543701, 43.110163243903585},
			{25.728735923767093, 43.10712416198819},
			{25.724401473999023, 43.10865938717618},
			{25.7244873046875, 43.11110313559475},
		}},
	)
	point := LatLng{Lat: 43.1089613, Lng: 25.7267396}
	for n := 0; n < b.N; n++ {
		polygon.Contains(point)
	}
}

func BenchmarkCountryMapPolygonContains(b *testing.B) {
	fc := readCountryGeoJson()
	polygon := NewPolygon(fc.Features[0].Geometry.Polygon)
	point := LatLng{Lat: 43.1089613, Lng: 25.7267396}
	for n := 0; n < b.N; n++ {
		polygon.Contains(point)
	}
}

func readCountryGeoJson() *geojson.FeatureCollection {
	f, err := os.Open("geo_testdata/country.geojson")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, _ := io.ReadAll(f)
	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		panic(err)
	}
	return fc
}
