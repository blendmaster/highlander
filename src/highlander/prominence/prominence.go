/* Detects prominent topology features from a 2d heightmap.
 *
 * Author: Steven Ruppert
 * For CSCI 447 Scientific Visualization, Spring 2013
 * Colorado School Of Mines
 */
package prominence

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"sort"
)

// I really don't get go's sorting API
type Pixels []*Pixel

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

// the sorted pixels (Descending) as well as a map from Locations to pixels
// from the image
func pixelsOf(heightmap image.Image) []*Pixel {
	b := heightmap.Bounds()
	size := b.Dx() * b.Dy()

	sorted := make([]*Pixel, size)

	i := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			// XXX explicit type assertion as color.Gray16, to be able to get
			// a number out of a color
			p := &Pixel{x, y, heightmap.At(x, y).(color.Gray16).Y}
			sorted[i] = p
			i++
		}
	}

	sort.Sort(Descending{sorted})
	return sorted
}

// greater height, with same position-based tiebreaker as pixel sorting to
// prevent strange things (from happening to me)
func gth(a *Feature, b *Feature) bool {
	if a.Height != b.Height {
		return a.Height > b.Height
	} else if a.X != b.X {
		return a.X < b.X
	}
	return a.Y < b.Y
}

// Print the prominent topologic features of a heightmap (as an Image).
// `threshold` controls which features will be returned.
func PrintProminentFeatures(heightmap image.Image, threshold uint16) {
	log.Println("sorting pixels")
	pixels := pixelsOf(heightmap)
	log.Println("pixels sorted!")

	aboveWater := make(map[Location]*Island)

	// since the list of features isn't stored but println'd when the prominence
	// is determined, output the absolute highest feature now, since its
	// prominence never changes
	fmt.Println(&Feature{pixels[0].X, pixels[0].Y, 65535, pixels[0].Height, Peak})

	for i, p := range pixels {
		if i%1000000 == 0 {
			log.Println("on pixel", i)
		}

		// lookup 4 connected pixels
		n, foundN := aboveWater[Location{p.X, p.Y - 1}]
		e, foundE := aboveWater[Location{p.X + 1, p.Y}]
		s, foundS := aboveWater[Location{p.X, p.Y + 1}]
		w, foundW := aboveWater[Location{p.X - 1, p.Y}]

		// set of connected islands
		connected := make(map[*Island]bool)

		if foundN {
			connected[Find(n)] = true
		}
		if foundE {
			connected[Find(e)] = true
		}
		if foundS {
			connected[Find(s)] = true
		}
		if foundW {
			connected[Find(w)] = true
		}

		switch len(connected) {
		case 0:
			// new Peak
			island := NewIsland()
			aboveWater[Location{p.X, p.Y}] = island
			island.HighestPeak = &Feature{p.X, p.Y, 65535, p.Height, Peak}
		case 1:
			// simple merge, loop only runs once
			for land := range connected {
				aboveWater[Location{p.X, p.Y}] = land
			}
		default:
			// 2 or more unconnected islands
			// saddle creation
			saddle := &Feature{p.X, p.Y, 65535, p.Height, Saddle}
			var highest, secondHighest *Feature
			var highestLand *Island

			for land := range connected {
				if secondHighest == nil || land.HighestPeak.Height > secondHighest.Height {
					secondHighest = land.HighestPeak
				}
				if highest == nil || land.HighestPeak.Height > highest.Height {
					highestLand = land
					secondHighest = highest
					highest = land.HighestPeak
				}
			}

			saddle.Prominence = secondHighest.Height - p.Height

			// update prominences, and union
			for land := range connected {
				if land.HighestPeak != highest {
					land.HighestPeak.Prominence = land.HighestPeak.Height - p.Height
					if land.HighestPeak.Prominence > threshold {
						fmt.Println(land.HighestPeak)
					}
					Union(land, highestLand)
				}
			}

			aboveWater[Location{p.X, p.Y}] = highestLand

			if saddle.Prominence > threshold {
				fmt.Println(saddle)
			}
		}
	}

	log.Println("All features found!")
}
