package traversal

import(
	"fmt"
	ON "github.com/fbaube/orderednodes"
)

// https://eli.thegreenplace.net/2023/preview-ranging-over-functions-in-go/

// Experimental package iter exports 
// type Seq[V any] func(yield func(V) bool)
// enabling (where f has type Seq[V], v has type V)
// for v := range f { ... } 

// DFS() should return an iterator
// that walks a tree in Preorder DFS.
/*
func counter() iter.Seq[T] [T ~ON.Nord] {
	// let's omit this for now
}

func main() {
	for n := range counter() {
		if n > 5 {
			break
		}
		fmt.Println(n)
	}
}
*/

func DFS(t *ON.SNord, doAndContinue func(*ON.SNord) bool) bool {
// func DFS[T ~ON.OrderedNode]doAndContinue func(*SNord) bool) bool {
  // Never executes ? Or maybe it can be
  // called this way by an iterator ?? 
  if t == nil {
    return true
  }
  if !t.HasKids() { return doAndContinue(t) }
  var p *ON.SNord 
  p = t.FirstKid().(*ON.SNord)
  for p != nil {
      fmt.Printf("p: %#v \n", *p)
      if false == DFS(p, doAndContinue) { return false }
      p = p.NextKid().(*ON.SNord)
      }
  return true 
}

/*
func (t *GON[E]) DFS(yield func(*GON[E]) bool) bool {
  // Never executes ? 
  if t == nil {
    return true
  }
  if !t.HasKids() { return yield(t) }
  var p *GON[E]
  p = t.FirstKid().(*GON[E])
  for p != nil {
      fmt.Printf("p: %#v \n", *p)
      if false == p.DFS(yield) { return false }
      p := p.NextKid().(*GON[E])
      }
  return true 
}
*/