// Disjoint set data structure as islands
package prominence

type Island struct {
  Pixel Pixel
  Parent *Island
}

func NewIsland() *Island {
  island := new(Island)
  island.Parent = island
  return island
}

func Find(island *Island) *Island {
  if island.Parent == island {
    return island
  }
  island.Parent = Find(island.Parent)
  return island.Parent
}

func Union(i1, i2 *Island) {
  root1, root2 := Find(i1), Find(i2)

  root1.Parent = root2
}
