// data structures
package prominence

type FeatureType int

const (
	Saddle FeatureType = iota
	Peak
)

// A topologic Saddle or Peak at a certain position.
type Feature struct {
	X, Y       int
	Prominence uint16
	Height     uint16
	Type       FeatureType
}

// Each height reading from the input image
type Pixel struct {
	X, Y   int
	Height uint16
}

// for lookup into a hash table
type Location struct{ X, Y int }
