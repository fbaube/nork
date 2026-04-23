https://eli.thegreenplace.net/2023/preview-ranging-over-functions-in-go/

type Tree[E any] struct {
  val         E
  left, right *Tree[E]
}

func (t *Tree[E]) Inorder(yield func(E) bool) bool {
  if t == nil {
    return true
  }
  return t.left.Inorder(yield) && yield(t.val) && t.right.Inorder(yield)
}

