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
  "log"
  "fmt"
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

// up to 2 unique pointers if a and b are unique
func uniqueIslands(a, b *Island) (w, x *Island) {
  if a == nil && b == nil {
    return nil, nil
  }
  if a != nil && b == nil {
    return a, nil
  }
  if a == nil && b != nil {
    return b, nil
  }
  if Find(a) == Find(b) {
    return a, nil
  }
  return a, b
}

// similar to above
func uniqueIslands3(a, b, c *Island) (_, _, _ *Island) {
  w, x := uniqueIslands(a, b)
  if w == nil { // a,b, are null
    return c, nil, nil
  }
  if x == nil { // w is not null
    y, z := uniqueIslands(w, c)
    return y, z, nil
  }
  // a, b are non, nil and unique
  if c == nil {
    return a, b, nil
  }
  rootc := Find(c)
  if rootc == Find(a) || rootc == Find(b) {
    return a, b, nil
  }
  return a, b, c
}

// same as above
func uniqueIslands4(a, b, c, d *Island) (_, _, _, _ *Island) {
  w, x, y := uniqueIslands3(a, b, c)
  if w == nil { // a, b, c are nil
    return d, nil, nil, nil
  }
  if x == nil { // w is not nil
    p, q := uniqueIslands(w, d)
    return p, q, nil, nil
  }
  if y == nil { // w, x are not nil
    p, q, r := uniqueIslands3(w, x, d)
    return p, q, r, nil
  }
  // a, b, c are all not nil
  if d == nil {
    return a, b, c, nil
  }
  rootd := Find(d)
  if rootd == Find(a) || rootd == Find(b) || rootd == Find(c) {
    return a, b, c, nil
  }

  // all are unique
  return a, b, c, d
}

// output a peak whose island joined at a saddle
func outputPeak(peak *Feature, saddleHeight, threshold uint16) {
  prominence := peak.Height - saddleHeight
  if prominence > threshold {
    fmt.Printf("%v, %v, %v, %v, %v\n",
      peak.X,
      peak.Y,
      prominence,
      peak.Height,
      Peak)
  }
}

func sort2(a, b *Island) (highest, secondHighest uint16) {
  if gth(a.HighestPeak, b.HighestPeak) {
    return a.HighestPeak.Height, b.HighestPeak.Height
  }
  // b is highest
  return b.HighestPeak.Height, a.HighestPeak.Height
}

func sort3(a, b, c *Island) (highest, secondHighest uint16) {
  abh, abs := sort2(a, b)

  cheight := c.HighestPeak.Height
  
  if cheight > abh {
    return cheight, abh
  }
  if cheight > abs {
    return abh, cheight
  }
  return abh, abs
}

func sort4(a, b, c, d *Island) (highest, secondHighest uint16) {
  abch, abcs := sort3(a, b, c)

  dheight := d.HighestPeak.Height
  if dheight > abch {
    return dheight, abch
  }
  if dheight > abcs {
    return abch, dheight
  }
  return abch, abcs
}

// The prominent topologic features of a heightmap (as an Image).
// `threshold` controls which features will be returned.
func PrintProminentFeatures(heightmap image.Image, threshold uint16) {
  log.Println("sorting pixels")
	pixels := pixelsOf(heightmap)
  log.Println("pixels sorted!")

	aboveWater := make(map[Location]*Island)

  // since the list of features isn't stored but println'd when the prominence
  // is determined, output the absolute highest feature now, since its
  // prominence never changes
  outputPeak(&Feature{pixels[0].X, pixels[0].Y, 65535, pixels[0].Height, Peak}, 0, threshold)

	for i, p := range pixels {
    if i % 1000000 == 0 {
      log.Println("on pixel", i)
    }

		// lookup 4 connected pixels
		n := aboveWater[Location{p.X, p.Y - 1}]
		e := aboveWater[Location{p.X + 1, p.Y}]
		s := aboveWater[Location{p.X, p.Y + 1}]
		w := aboveWater[Location{p.X - 1, p.Y}]

    a, b, c, d := uniqueIslands4(n, e, s, w)

		switch {
		case a == nil:
			// new Peak
      island := NewIsland()
      aboveWater[Location{p.X, p.Y}] = island
			island.HighestPeak = &Feature{p.X, p.Y, 65535, p.Height, Peak}
		case b == nil:
			// simple merge, loop only runs once
      aboveWater[Location{p.X, p.Y}] = a
    case c == nil:
      //2 connected islands, a and b
      highest, secondHighest := sort2(a, b)

      if a.HighestPeak.Height == highest {
        outputPeak(b.HighestPeak, p.Height, threshold)
        Union(b, a)
        aboveWater[Location{p.X, p.Y}] = a
      } else {
        outputPeak(a.HighestPeak, p.Height, threshold)
        Union(a, b)
        aboveWater[Location{p.X, p.Y}] = b
      }

      prominence := secondHighest - p.Height
      // output saddle
      if prominence > threshold {
        fmt.Printf("%v, %v, %v, %v, %v\n", p.X, p.Y, prominence, p.Height, Saddle)
      }
    case d == nil:
      //3 connected islands, a and b and c
      highest, secondHighest := sort3(a, b, c)

      switch (highest) {
      case a.HighestPeak.Height:
        outputPeak(b.HighestPeak, p.Height, threshold)
        outputPeak(c.HighestPeak, p.Height, threshold)
        Union(b, a)
        Union(c, a)
        aboveWater[Location{p.X, p.Y}] = a
      case b.HighestPeak.Height:
        outputPeak(a.HighestPeak, p.Height, threshold)
        outputPeak(c.HighestPeak, p.Height, threshold)
        Union(a, b)
        Union(c, b)
        aboveWater[Location{p.X, p.Y}] = b
      case c.HighestPeak.Height:
        outputPeak(a.HighestPeak, p.Height, threshold)
        outputPeak(b.HighestPeak, p.Height, threshold)
        Union(a, c)
        Union(b, c)
        aboveWater[Location{p.X, p.Y}] = c
      }

      prominence := secondHighest - p.Height
      // output saddle
      if prominence > threshold {
        fmt.Printf("%v, %v, %v, %v, %v\n", p.X, p.Y, prominence, p.Height, Saddle)
      }
    default:
      log.Println("on pixel", p)
      log.Println("saddle", a, b, c, d)
      // a, b, c, d all connected
      highest, secondHighest := sort4(a, b, c, d)

      switch (highest) {
      case a.HighestPeak.Height:
        outputPeak(b.HighestPeak, p.Height, threshold)
        outputPeak(c.HighestPeak, p.Height, threshold)
        outputPeak(d.HighestPeak, p.Height, threshold)
        Union(b, a)
        Union(c, a)
        Union(d, a)
        aboveWater[Location{p.X, p.Y}] = a
      case b.HighestPeak.Height:
        outputPeak(a.HighestPeak, p.Height, threshold)
        outputPeak(c.HighestPeak, p.Height, threshold)
        outputPeak(d.HighestPeak, p.Height, threshold)
        Union(a, b)
        Union(c, b)
        Union(d, b)
        aboveWater[Location{p.X, p.Y}] = b
      case c.HighestPeak.Height:
        outputPeak(a.HighestPeak, p.Height, threshold)
        outputPeak(b.HighestPeak, p.Height, threshold)
        outputPeak(d.HighestPeak, p.Height, threshold)
        Union(a, c)
        Union(b, c)
        Union(d, c)
        aboveWater[Location{p.X, p.Y}] = c
      case d.HighestPeak.Height:
        outputPeak(a.HighestPeak, p.Height, threshold)
        outputPeak(b.HighestPeak, p.Height, threshold)
        outputPeak(c.HighestPeak, p.Height, threshold)
        Union(a, d)
        Union(b, d)
        Union(c, d)
        aboveWater[Location{p.X, p.Y}] = d
      }

      prominence := secondHighest - p.Height
      // output saddle
      if prominence > threshold {
        fmt.Printf("%v, %v, %v, %v, %v\n", p.X, p.Y, prominence, p.Height, Saddle)
      }
		}
	}

  log.Println("All features found!")
}
