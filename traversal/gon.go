package traversal

import(
	ON "github.com/fbaube/orderednodes"
	)

// GON is a Generic-enabled [Nord] (ordered node).
type GON [T any] struct {
     ON.Nord
     // This field needs to have a name, to avoid the compiler error
     // "embedded field type cannot be a (pointer to a) type parameter".
     // Use "V", for "variable" or "variant" (or "vary-er"). 
     V T
     }

type SON GON[string]

// type GenON ~GON

// NewSON creates a new simple String-only
// Ordered Node, and returns a pointer to it. 
func NewSON(s string) *SON { // *GON[string] {
     p := new(SON) // GON[string])
     p.V = s
     return p
     }

