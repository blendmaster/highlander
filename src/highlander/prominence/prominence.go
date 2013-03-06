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
	"sort"
  "container/list"
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

// the sorted pixels (Descending) as well as a map from Locations to pixels
// from the image
func pixelsOf(heightmap image.Image) []Pixel {
	b := heightmap.Bounds()
	size := b.Dx() * b.Dy()

	sorted := make([]Pixel, size)

	i := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			// XXX explicit type assertion as color.Gray16, to be able to get
			// a number out of a color
			p := Pixel{x, y, heightmap.At(x, y).(color.Gray16).Y}
			sorted[i] = p
			i++
		}
	}

	sort.Sort(Descending{sorted})
	return sorted
}

// The prominent topologic features of a heightmap (as an Image).
// `threshold` controls which features will be returned.
func ProminentFeatures(heightmap image.Image, threshold uint16) *list.List {
	_ = fmt.Println

	pixels := pixelsOf(heightmap)

	aboveWater := make(map[Location]*Island)

	features := list.New()

	for _, p := range pixels {

		island := NewIsland()
		aboveWater[Location{p.X, p.Y}] = island

		// lookup 4 connected pixels
		n, foundN := aboveWater[Location{p.X, p.Y - 1}]
		e, foundE := aboveWater[Location{p.X + 1, p.Y}]
		s, foundS := aboveWater[Location{p.X, p.Y + 1}]
		w, foundW := aboveWater[Location{p.X - 1, p.Y}]

    // FIXME two or more of the connected pixels could be from the same island,
    // need to update logic for that.
		numFound := 0
		if foundN {
			numFound++
		}
		if foundE {
			numFound++
		}
		if foundS {
			numFound++
		}
		if foundW {
			numFound++
		}

		if numFound == 0 {
			// new Peak
      island.HighestPeak = &Feature{p.X, p.Y, 65535, p.Height, Peak}
			features.PushBack(island.HighestPeak)
		} else if numFound == 1 {
			// simple merge
			switch {
			case foundN:
				Union(island, n)
        island.HighestPeak = n.HighestPeak
			case foundE:
				Union(island, e)
        island.HighestPeak = e.HighestPeak
			case foundS:
				Union(island, s)
        island.HighestPeak = s.HighestPeak
			case foundW:
				Union(island, w)
        island.HighestPeak = w.HighestPeak
			}
		} else {
			// saddle creation
			saddle := &Feature{p.X, p.Y, 65535, p.Height, Saddle}
			var highest, secondHighest, hNE, sNE, hSW, sSW *Feature

			switch {
			case foundN || foundE:
				if foundN && foundE {
					if n.HighestPeak.Height > e.HighestPeak.Height {
						hNE = n.HighestPeak
						sNE = e.HighestPeak
					} else {
						hNE = e.HighestPeak
						sNE = n.HighestPeak
					}
          Union(e, n)
				} else {
          if foundN {
            hNE = n.HighestPeak
          } else if foundE {
            hNE = e.HighestPeak
          }
				}
			case foundS || foundW:
				if foundS && foundW {
					if s.HighestPeak.Height > w.HighestPeak.Height {
						hSW = s.HighestPeak
						sSW = w.HighestPeak
					} else {
						hSW = w.HighestPeak
						sSW = s.HighestPeak
					}
          Union(w, s)
				} else {
          if foundS {
            hSW = s.HighestPeak
          } else if foundW {
            hSW = w.HighestPeak
          }
				}
			}

			switch {
			case hNE != nil && hSW != nil:
				if hNE.Height > hSW.Height {
          highest = hNE
					secondHighest = hSW
				} else {
          highest = hSW
					secondHighest = hNE
				}
				if sNE != nil && sNE.Height > secondHighest.Height {
					secondHighest = sNE
				}
				if sSW !=nil && sSW.Height > secondHighest.Height {
					secondHighest = sSW
				}

        // merge islands, second level
        switch {
        case foundN && foundS:
          Union(s, n)
        case foundN && foundW:
          Union(w, n)
        case foundE && foundS:
          Union(s, e)
        case foundE && foundW:
          Union(w, e)
        }

			case sNE != nil && hNE != nil && hSW == nil:
        highest = hNE
        secondHighest = sNE
			case sSW != nil && hNE == nil && hSW != nil:
        highest = hSW
        secondHighest = sSW
      default: 
        switch {
        case foundN && foundS:
          if n.HighestPeak.Height > s.HighestPeak.Height {
            highest = n.HighestPeak
            secondHighest = s.HighestPeak
          } else {
            highest = s.HighestPeak
            secondHighest = n.HighestPeak
          }
        case foundN && foundW:
          if n.HighestPeak.Height > w.HighestPeak.Height {
            highest = n.HighestPeak
            secondHighest = w.HighestPeak
          } else {
            highest = w.HighestPeak
            secondHighest = n.HighestPeak
          }
        case foundE && foundS:
          if e.HighestPeak.Height > s.HighestPeak.Height {
            highest = e.HighestPeak
            secondHighest = s.HighestPeak
          } else {
            highest = s.HighestPeak
            secondHighest = e.HighestPeak
          }
        case foundE && foundW: // must be true
          if e.HighestPeak.Height > w.HighestPeak.Height {
            highest = e.HighestPeak
            secondHighest = w.HighestPeak
          } else {
            highest = w.HighestPeak
            secondHighest = e.HighestPeak
          }
        }
			}

      fmt.Println("secondHighest", secondHighest.Height, "p", p.Height, secondHighest.Height - p.Height)
			saddle.Prominence = secondHighest.Height - p.Height

      // update prominences, then highest peaks
      if foundN && highest != n.HighestPeak {
        n.HighestPeak.Prominence = n.HighestPeak.Height - p.Height
        n.HighestPeak = highest
      }
      if foundE && highest != e.HighestPeak {
        e.HighestPeak.Prominence = e.HighestPeak.Height - p.Height
        e.HighestPeak = highest
      }
      if foundS && highest != s.HighestPeak {
        s.HighestPeak.Prominence = s.HighestPeak.Height - p.Height
        s.HighestPeak = highest
      }
      if foundW && highest != w.HighestPeak {
        w.HighestPeak.Prominence = w.HighestPeak.Height - p.Height
        w.HighestPeak = highest
      }

      // _this_ pixel 
      island.HighestPeak = highest

      features.PushBack(saddle)
		}
	}

  thresholded := list.New()
  // filter out features not meeting the threshold
	for e := features.Front(); e != nil; e = e.Next() {
    if e.Value.(*Feature).Prominence > threshold {
      thresholded.PushBack(e.Value)
    }
	}

	return thresholded
}
