/* Detects prominent topology features from a 2d heightmap.
 *
 * Author: Steven Ruppert
 * For CSCI 447 Scientific Visualization, Spring 2013
 * Colorado School Of Mines
 */
package prominence

import (
	"image"
	"image/color"
	"sort"
)

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

// I really don't get go's sorting API
type Pixels []Pixel

func (p Pixels) Len() int      { return len(p) }
func (p Pixels) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type Descending struct{ Pixels }

// reverse sort (decending) so using greater than
// but sort earlier pixels (by positon) before later pixels
func (p Descending) Less(i, j int) bool {
	a := p.Pixels[i]
	b := p.Pixels[j]
	if a.Height != b.Height {
		return a.Height > b.Height
	} else if a.X != b.X {
		return a.X < b.X
	}
	return a.Y < b.Y
}

// In decreasing order, the pixels from the Gray16 Image.
func sortedPixels(heightmap image.Image) []Pixel {
	b := heightmap.Bounds()
	size := b.Dx() * b.Dy()

	sorted := make([]Pixel, size)

	i := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			sorted[i] = Pixel{x, y, heightmap.At(x, y).(color.Gray16).Y}
			i++
		}
	}

	sort.Sort(Descending{sorted})
	return sorted
}

// The prominent topologic features of a heightmap (as an Image).
// `threshold` controls which features will be returned.
func ProminentFeatures(heightmap image.Image, threshold int) []Feature {

	return []Feature{{0, 0, 0, 0, Peak}}
}
