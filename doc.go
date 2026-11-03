// Package nork is Nodes with ORdered Kids - a building
// block for hierarchical trees of nodes, where sibling 
// nodes are ordered and keep their order.
//
// A Nork stores bidirectional parent/kid relationships, 
// plus other useful hierarchy-related info like levels 
// and paths. A Nork does not need or use Go generics.
//
// Node order is of course important for markup languages 
// in general and XML mixed content in particular, while 
// unimportant when XML is used purely for data records.
//
// Child nodes ("Kid"s) have an externally-defined order,
// implying that child nodes are not directly [Comparable] 
// amongst themselves. This specific externally-provided 
// order of child nodes is essential for representing content.
//
// (Child ordering can help in other contexts too, such
// as filesystem operations, but it turns out that the
// Go stdlib generally returns directory items in an 
// order - lexical order - and walks directories in
// lexical order.)
//
// The ordering lets us define funcs like FirstKid, NextPeer,
// PrevPeer, LastKid. They are defined in interface [Norker],
// which is implemented not by struct [Nork] but rather by the
// pointer [*Nork]. This makes nodes writable and also sharable.
//
// Advanced Usage: Directions, Use Cases, Func Aliasing 
// 
// Nork can become a "QuadraNode", handling four directions 
// of connections for four distinct node user cases (UCs). 
// Nork will at first be kind of a maximalist implementation, 
// subject later to redefinition and getting slimmed down.
//
// Paths need not be in the filesystem, and so paths are
// IAW package [path] rather than package [path/filepath].
// 
// Structs that embed struct Nork can use funcs to alias
// the names and usage of fields in Nork. For example, a 
// Nork embedded in a [Cnty] (Content Entity) used for a
// file-or-dir can provide access to non-public field
// `relPath` via `func RelFP()` (where FP is "filepath"). 
//
// Quadranode Discussion
// 
// Directions in GUI: 
//  1. Left  is ParentDir and Incoming Links
//  2. Right is Kids and Outgoing Links
//  3. Above is Incomiong Facets and maybe other metadata
//  4. Below is Versions 
//
// Fields and use cases: 
//  - The field Prnt is for a "Parent" singleton of some type (whenever
//    applicable).
//  - The fields Tags and Vers are used pretty much the same in all UC's.
//  - UC.FSO: File System Object: represent a dir-or-file-or-softlink:
//    	      Here ordering is less important. Paths are valid FS paths. 
//            Prnt:dir, Kids:contents(!dir-v-file), Usrs:incoming-symlinks 
//  - UC.CMS: CMS Usage: Table of Contents line item (transclusion in CMS)
//    	      This should be an ideal use case. 
//            Prnt:?TBD, Kids:sub-ToC's/outrefs, Usrs:inrefs+transcluders
//  - UC.XML: XML text: tag/element in markup (XML, or other w AST)
//    	      This has complexity handling same-named siblings, such as
//	      multiple <p> tags. Paths are tricky (mult. same tags).
//            Prnt:parent-elm(or file @root), Kids:kid-elms, Usrs:?entity-stuf
//  - UC.GST: GolangAST: Go code AST node (https://pkg.go.dev/go/ast#Node)
//          Pnt:node(or root's file), Kids:AST, Usrs:callers/refs
//
// Memory and Nodes (NOTE: Obsolete discussion) 
// 
// There are two distinct memory management nodes for allocating and
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
// Representation and Storage of Links Between Nodes
// 
// There are multiple ways to represent node trees in our
// SQLite DBMS, and multiple ways to walk a node tree, so
// there is a unavoidable complexity wherever we look.
//
// NOTE: DOM markup exhibits name duplication: In UC.1 we never have two
// same-named files in the same directory, but in UC.2 we might have (for
// example) multiple sibling <p> tags. So when representing markup, a map
// from paths to Norks would fail unless the tags are made unique with
// subscript indices (such as "[1]", "[2]").
//
// Link fields are lower-cased so that other packages cannot damage links.
//
// NOTE: This implementation stores pointers to child nodes in a doubly
// linked list, not a slice. Therefore a Nork does not have a complete
// set of pointers to all of its kids. Therefore (a) it is not simple
// to get a kid count, because it requires a list traversal, and
// (b) it is not feasible to modify this code to define a simpler,
// more efficient variant of Nork that has unordered kids.
//
// UUID version 7
// 
// UUID (version 7, featuring time ordering) has outsized importance,
// because it can be used to reference all manner of other DB tables
// and data structures and external (to the system) targets.
//
// This could work nicely with the Quadranode concept.
//
// A UUID should be absolutely unique within a particular database,
// and quite possible between databases too.
//
// Therefore it should be used as the main indexing field for and
// and all objects in SQLite.
//
// It will be necessary to add prefixes to UUIDs to specify (e.g.) 
// the DB tables of targets.
//
// This requires the use of a custom UUID v7 implementation: WIP!
//
// The version defined here replaces much of a UUID's random with
// fields that embed the item's database table name (or short name)
// and row ID. The former will add great efficienty; the latter may
// turn out to be unnecessary and unhelpful.
//
// When a tree of items is loaded into memory, a UUID is a very
// inefficient way to follow a link between nodes, so instead,
// there should be pointers to other nodes, which (of course)
// are session-specific and should not be persisted to database. 
//
// It is straightforward to store 1-to-1 links in a database.
// But there won't be many of those, if any. But there will be
// many cases of a parent with many kids, and for a Quadranode
// there will be many nodes with multiple parents. (Altho we
// do not expect to see N-to-N situations.)
//
// So for these two situations - 1-to-N, and in the future
// N-to-1 - there can be two database tables "Links1toN" 
// and later "LinksNto1". And because every UUID7 will
// contain a table name corresponding to an item type,
// use cases where the items are a mix of types can
// simply store them all together. 
// .
package nork
