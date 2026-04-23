// https://eli.thegreenplace.net/2023/preview-ranging-over-functions-in-go/

type Tree[E any] struct {
  val         E
  kids []*Tree[E]
}

func (t *Tree[E]) Inorder(yield func(E) bool) bool {
  if t == nil {
    return true
  }
  return t.left.Inorder(yield) && yield(t.val) && t.right.Inorder(yield)
}

func (t *Tree[E]) Topdown(yield func(E) bool) bool {
  if t == nil {
    return true
  }
  return yield(t.val) && t.left.Preorder(yield) && t.right.Preorder(yield)
}

func (t *Tree[E]) Postorder(yield func(E) bool) bool {
  if t == nil {
    return true
  }
  return t.left.Postorder(yield) && t.right.Postorder(yield) && yield(t.val) 
}

