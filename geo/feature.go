package geo

type Feature struct {
	Type       string
	Geometry   Geometry
	Properties *Properties
}
