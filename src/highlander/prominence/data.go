// data structures
package prominence

import "fmt"

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

// output feature as a (lite) csv
func (f *Feature) String() string {
  return fmt.Sprintf("%v, %v, %v, %v, %v", f.X, f.Y, f.Prominence, f.Height, f.Type)
}


// Each height reading from the input image
type Pixel struct {
	X, Y   int
	Height uint16
}

// for lookup into a hash table
type Location struct{ X, Y int }
