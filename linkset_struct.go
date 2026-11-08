package nork

import(
	FU "github.com/fbaube/fileutils"
)

// Linkset stores the links for a node. It is designed to
// work both on disk (using UUIDs) and in-memory (using
// pointers). When persisting to disk, the fields ( prnt
// prev next ) are sufficient to reconstruct a tree that
// has all kid nodes in their correct order.
//
// There needs to be a distinguished "empty value" to
// mark a non-value, so for a pointer this should be Go
// nil and for a UUID it can be the zero-length string.
//
// When a tree structure is built that extends upward,
// with the root at the bottom, maybe some of these
// field names should be changed to be less confusing. 
// .
type Linkset[T comparable] struct {
   prnt   T  // left: optional: shown at top of L-H list
   prev   T  // previous sibling, or nil if is the first 
   next   T  // next sibling,  or nil if it is the last 
   kids []T  // rite: outgoing: many, in an order 
   usrs []T  // left: incoming: symlinks(FSO)/referrers(CMS)
   tags []T  // topp: incoming: tags/facets
   vers []T  // botm: versions: is a list, not a tree 
}

// Null returns the zero value for any T.
// https://www.perplexity.ai/search/717b5c69-194e-4599-88d6-2f00ee7a236a 
func Null[T comparable]() T {
     // zed is not initialized to any specific valur,
     // so it gets initialized to the "zero value" !
     var zed T
     return zed
}

// IsNull reports whether arg is the zero value for T.
// https://www.perplexity.ai/search/717b5c69-194e-4599-88d6-2f00ee7a236a 
// NOTE: "==" works byte-for-byte for UUIDs !!
// See https://www.perplexity.ai/search/717b5c69-194e-4599-88d6-2f00ee7a236a 

func IsNull[T comparable](arg T) bool {
    var zed T
    return arg == zed
}

