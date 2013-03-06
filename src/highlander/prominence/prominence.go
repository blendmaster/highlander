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
  "fmt"
)

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

// In decreasing order by height, the pixels from the Gray16 Image.
func sortedPixels(heightmap image.Image) []Pixel {
	b := heightmap.Bounds()
	size := b.Dx() * b.Dy()

	sorted := make([]Pixel, size)

	i := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
      // XXX explicit type assertion as color.Gray16, to be able to get
      // a number out of a color
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
  _ = fmt.Println
  pixels := sortedPixels(heightmap)
  features := make([]Feature, len(pixels))
  for i, p := range(pixels) {
    // TODO stuff
    features[i] = Feature{p.X, p.Y, 0, p.Height, Peak}
  }

	return features
}
