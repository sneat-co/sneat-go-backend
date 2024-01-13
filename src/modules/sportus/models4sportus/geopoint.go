package models4sportus

// Degrees are units of angular measure.
type Degrees float64

// GeoPoint represents GPS coordinates.
type GeoPoint struct {
	Latitude  Degrees `json:"latitude" firestore:"lat"`
	Longitude Degrees `json:"longitude" firestore:"long"`
}
