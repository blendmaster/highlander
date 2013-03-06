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
	Prominence int
	Height     uint16
	Type       FeatureType
}

// Each height reading from the input image
type Pixel struct {
	X, Y   int
	Height uint16
}
