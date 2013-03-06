/* Detects prominent topology features from a 2d heightmap.
 *
 * Author: Steven Ruppert
 * For CSCI 447 Scientific Visualization, Spring 2013
 * Colorado School Of Mines
 */
package prominence

import "image"

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

// The prominent topologic features of a heightmap (as an Image).
// `threshold` controls which features will be returned.
func ProminentFeatures(heightmap image.Image, threshold int) []Feature {
	return []Feature{{0, 0, 0, 0, Peak}}
}
