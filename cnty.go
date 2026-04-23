// nork is a Generic Node, Ordered Kids. It stores
// basic bidirectional parent/kid relationships,
// plus other useful hierarchy-related info.
//
// interface [Norker] is implemented not for struct
// [Nork] but rather for the pointer, i.e. `*Nork´. 
// This makes nodes writable and also sharable.
//
package nork

import(
	"fmt"
	"os"
)

// Nord is shorthand for "ordered node" (or "ordinal node") -
// a Node whose child nodes have an externally-defined order,
// so that child nodes are not directly [Comparable] amongst
// themselves. This specific specified order or child nodes 
// is essential for representation of content.
//
// (Child ordering can help in other contexts too, such 
// as filesystem operations, but it turns out that the 
// Go stdlib generally returns directory items in lexical 
// order and walks directories in lexical order.)
// 
// The ordering lets us define funcs like FirstKid, NextKid, 
// PrevPeer, LastKid. They are defined in interface [Norker].
// A Nork for file/dir also contains its own relative path
// (relative to its inbatch) and absolute path. 
//
// Usage 
//
// Methods on Nork need to be methods on pointers rather
// than methods on values. But the interface [Norker] 
// will be defined on pointers rather than values, so
// the resulting syntax will be visually acceptable. 
//
// Alternatives Everywhere
// 
// There are three use cases identified for Norks:
//  - UC.1: files & directories (here ordering is less important) 
//  - UC.2: XML/HTML/DOM markup (this has complexity handling
//          same-named siblings, such as multiple <p> tags)
//  - UC.3: [Lw]DITA map files and other ToC's (these should
//          be an ideal use case) 
//
// Also there are two distinct memory management nodes for allocating and
// linking nodes:
//  - The "traditional method" of allocating nodes individually, and linking
//    them using pointers. Using this method, both deletions and insertions
//    are relatively simple.
//     - Such nodes can also be loaded into a map, for random access to nodes
//       based on path. 
//  - The "new-fangled way" called an "arena", where we put all our nodes in
//    a big slice, and link them using indices. This method is much kinder
//    on memory management, but might becomes clumsy when we need dynamic
//    node management.
//     - Deletions are easy if we just zero out the slice entry; we cannot
//       then do compaction because it would require updating all indices
//       past the first point of compaction.
//     - Insertions are costly. However note that in this implementation,
//       for a node's kids, we use a linked list rather than a slice, so
//       this makes it easier to append a new node at the end of the
//       arena-slice and then update indices, wherever they may be
//       elsewhere in the slice.
//     - In any case, if an arena-slice has to grow (because of a call 
//       to append), it might be moved elsewhere in memory, which would
//       invalidate all ptrs to other Norks! If this happens, we should
//       trust only the "traditional method".
//
// Also there are multiple ways to represent node trees in our SQLite DBMS,
// and multiple ways to walk a node tree, so there is a unavoidable complexity
// wherever we look. 
// 
// NOTE: DOM markup exhibits name duplication: In UC.1 we never have two
// same-named files in the same directory, but in UC.2 we might have (for
// example) multiple sibling <p> tags. So when representing markup, a map
// from paths to Norks would fail unless the tags are made unique with 
// subscript indices (such as "[1]", "[2]").
//
// NOTE: Using Norks for files & dirs exhibits strong typing. Dirs are
// dirs and files are files and never the twain shall meet. This means
// that (a) dirs cannot contain own-content, and (b) files can never 
// be non-leaf nodes. (Note tho that symlinks have aspects of both.) 
// However this dir/file/etc typing is too complex to handle here in 
// a Nork, because a leaf node can be either a file or a dir, and a
// field like "canKid bool" is a bit OTT, so the file/dir distinction
// is handled instead by an outer struct type that embeds Nork, such 
// as [fileutils.FSItem]. 
//
// If we build up a tree of Norks when processing an [os.DirFS], the 
// strict ordering provided by DirFS is not strictly needed, BUT it 
// can anyways be used (and relied upon) because the three flavors of
// WalkDir are deterministic (using lexical order). WalkDir does tho
// promise that a given Nork will always appear AFTER the Nork for 
// its directory has appeared, which makes it "easy" to build a tree.
//
// Link fields are lower-cased so that other packages cannot damage links. 
//
// NOTE: This implementation stores pointers to child nodes in a doubly
// linked list, not a slice. Therefore a Nork does not have a complete 
// set of pointers to all of its kids. Therefore
//  1. it is not simple to get a kid count, because it requires a 
//     list traversal, and
//  2. it is not feasible to modify this code to define a simpler, 
//     more efficient variant of Nork that has unordered kids. 
// .
type Cnty struct {
     	// Nork is embedded as an instance, not as a pointer. 
     	Nork
     	// ====================================
	//  Substructures for: Adjacency Lists 
	// ====================================
	// Every implementation of adjacency lists points 
	// (1) to parent, (2) to peers in own list, and 
	// (3) to first & last kids in linked list of all its kids. 
     	// ------------------------------------------------
	//  Substructure for Adjacency List #1 based on Go
	//  ptrs-to-structs (not indices). Note that these 
	//  are provided IN PARALLEL with the ptr-to-struct 
	//  slices in the embedded Nork. This gives us
	//  redundancy for error checking, and maybe also
	//  lets us run performance comparisons. 
     	// ------------------------------------------------
//	parent             Norker // level up   // dupe of Nork.prnt 
	prevPeer, nextPeer Norker // level same // like Nork.kids []*Nork
	firstKid, lastKid  Norker // level down // like Nork.kids []*Nork
/*
     	// ------------------------------------------
	//  Substructure for Adjacency List #2 based
	//  on discrete indices into "arena" slice 
     	// ------------------------------------------
	iParent              int // level up
	iPrevPeer, iNextPeer int // level same 
	iFirstKid, iLastKid  int // level down
*/
/*
     	// ------------------------------------------
	//  Substructure for Adjacency List of KIDS
	//  where they are stored on a single string 
	//  field that holds all applicable indices.
	// ------------------------------------------
	// kidIdxs when empty is "," (or ""), else
	// e.g. ",1,4,56,". The kidIdxs should be in
	// the same order as the Kid nodes themselves.
	// Comma-bracketing simplifies search (",%d,").
	// ------------------------------------------
	kidIdxs string
*/
/*
     	// ----------------------------------
	//  Temporarily unused: three fields	
     	// ----------------------------------
	// seqID is a unique ID under this node's tree's root. It does not 
	// need to be the same as (say) the index of this Nork in a slice 
	// of Nork's, but it probably is. Its use is optional, and also 
	// it can be used in other ways in structs that embed Nork.
	// >> seqID int
	// parSeqID and kidSeqID's can add a layer of error checking 
	// and simplified access. Their use is optional.
	// kidSeqIds when empty is ",", otherwise e.g. ",1,4,56,". 
	// the seqIds should (well, can) be in the same order as 
	// the Kid nodes themselves. The bracketing by commas makes 
	// searching simpler (",%d,").
	// >> parSeqID, kidSeqID string
*/
}

// RootNork is defined, so that assignments
// to/from a root node have to be explicit.
// type RootNork[D any] Nork[D any]

// Root walks the tree upward until [IsRoot] is true,
// so it does not use any global variables.
func (p *Nork) Root() *Nork {
	if p.IsRoot() {
		return p
	}
	var ondr *Nork
	ondr = p
	for !ondr.IsRoot() {
		ondr = ondr.Parent()
	}
	return ondr
}

// setlevel is duh.
func (p *Nork) setLevel(i int) {
	p.level = i
}

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func (pOld *Nork) ReplaceBy(pNew *Nork) *Nork {
	// REPLACE SIBLINGS' SIBBLE-LINKS
	// REPLACE KIDS' PARENT-LINKS
	// REPLACE PARENT'S KID-LINK

	// We require that pNew has no existing links
	if pNew.PrevPeer() != nil || pNew.NextPeer() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in ReplaceBy: " +
			"Tag<< %+v >> new<< %+v >>\n", pOld, pNew)
		panic("ReplaceBy(K) can't cos K has siblings")
	}
	if pNew.Parent() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in ReplaceBy: " +
			"Tag<< %+v >> new<< %+v >>\n", pOld, pNew)
		panic("E.ReplaceBy(K) can't cos K has non-P parent")
	}
	// REPLACE SIBLINGS' SIBBLE-LINKS
	prv := pOld.PrevPeer()
	if prv != nil {
		pNew.SetPrevPeer(prv)
		prv.SetNextPeer(pNew)
	}
	nxt := pOld.NextPeer()
	if nxt != nil {
		pNew.SetNextPeer(nxt)
		nxt.SetPrevPeer(pNew)
	}
	// REPLACE KIDS' PARENT-LINKS
	if pOld.FirstKid() != nil {
		crntKid := pOld.FirstKid()
		for crntKid != nil {
			crntKid.SetParent(pNew)
			pNew.AddKid(crntKid)
			crntKid = crntKid.NextPeer()
		}
	}
	// REPLACE PARENT'S KID-LINK

	return pNew
}

/*
type  FSData struct { Path string; Size int64 }
type CMSData struct { Markdown string }
type XMLData struct { Tag string; Attrs map[string]string }
type GSTData struct { Kind string; Value string }
*/

