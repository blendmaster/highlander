package prominence

import (
	"sort"
	"testing"
)

func TestPixelSort(t *testing.T) {
	pixels := []Pixel{
		{0, 1, 1},
		{0, 0, 0},
		{0, 10, 0},
		{0, 1, 2},
		{1, 3, 6},
	}

	expected := []Pixel{
		{1, 3, 6},
		{0, 1, 2},
		{0, 1, 1},
		{0, 0, 0},
		{0, 10, 0},
	}

	sort.Sort(Descending{pixels})

	for i := range pixels {
		if pixels[i] != expected[i] {
			t.Errorf("not sorted: %v", pixels)
			break
		}
	}
}
