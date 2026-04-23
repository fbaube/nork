// nork is a Generics Node with Ordered Kids. It stores 
// basic bidirectional parent/kid relationships, plus other 
// useful hierarchy-related info like levels and paths. 
//
// interface [Norker] is implemented not for struct
// [Nork] but rather for the pointer, i.e. `*Nork´. 
// This makes nodes writable and also sharable.
//
package nork

import(
)

// StringFunc is used by interface Norker, so a 
// method signature actually (MAYBE!) looks like:
//   func (*Nork) FuncName() string
type StringFunc func(Norker) string

// Nork is a "QuadraNode", handling four directions of connections 
// for four distinct node user cases (UCs). The associated (generic) 
// data structure `D` is a pointer, so that the data stored in a
// Q4Node can be written-to and can be passed around (shared).
//
// QNode is kind of a maximalist implementation of nork,
// and is subject to redefinition and getting slimmed down.
//
// Child nodes ("Kid"s) have an externally-defined order,
// implying that child nodes are not directly [Comparable] 
// amongst themselves. This specific specified order or 
// child nodes is essential for representation of content.
//
// (Child ordering can help in other contexts too, such
// as filesystem operations, but it turns out that the
// Go stdlib generally returns directory items in an 
// order - lexical order - and walks directories in
// lexical order.)
//
// The ordering lets us define funcs like FirstKid, NextPeer,
// PrevPeer, LastKid. They are defined in interface [Norker].
//
// Structs that embed Nork can use funcs to redefine 
// the names and usage of fields in Nork. For example,
// a Nork embedded in a [Cnty] (Content Entity) used for
// a file-or-dir can provide access to non-public field
// `relPath` via `func RelFP()` (where FP is filepath). 
//
// Implementation of Interface 
//
// Methods on Nork need to be methods on pointers rather
// than methods on values. But the interface [Norker]
// will be defined on pointers rather than values, so
// the resulting syntax will be visually acceptable.
//
//
// When a Nork is embedded in another struct, the embedding struct
// can rename the functions defined here. For example, [Nork.relPath]
// can be renamed as ´relFP`when we are dealing with a filesystem.
//
// Directions in GUI: 
//  1. Left is ParentDir and Incoming Links
//  2. Right is Kids and Outgoing Links
//  3. Up is Incomiong Facets and maybe other metadata
//  4. Down is Versions 
//
// Fields and use cases: 
//  - The field Prnt is for a "Parent" singleton of some type (whenever
//    applicable).
//  - The fields Tags and Vers are used pretty much the same in all UC's.
//  - UC.F: File System Item: represent a dir-or-file-or-softlink:
//    	    Here ordering is less important.
//          Prnt:dir, Kids:contents(!dir-v-file), Usrs:incoming-symlinks 
//  - UC.C: CMS Usage: Table of Contents line item (transclusion in CMS)
//    	    This should be an ideal use case. 
//          Prnt:?TBD, Kids:sub-ToC's/outgoing-refs, Usrs:inrefs+transcluders
//  - UC.X: XML text: tag/element in markup (XML, or other w AST)
//    	    This has complexity handling same-named siblings, such as
//	    multiple <p> tags. 
//          Prnt:parent-elm(or root's file), Kids:kid-elms, Usrs:?entity-stuff
//  - UC.G: GolangAST: Go code AST node (https://pkg.go.dev/go/ast#Node)
//          Pnt:node(or root's file), Kids:AST, Usrs:callers/refs
//
// Note that using `*Nork[D]` everywhere might not work out,
// because the fields may end up being quite different from each other.
//
// Note that UUID might take on outsized importance, because it might be
// used to reference all manner of other DB tables and data structures.
// It may become necessary to add prefixes to UUIDs to specify the DB
// tables of targets.
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
// set of pointers to all of its kids. Therefore (a) it is not simple
// to get a kid count, because it requires a list traversal, and
// (b) it is not feasible to monify this code to define a simpler,
// more efficient variant of Nork that has unordered kids.
// .
type Nork struct {
//  ==============
//    BASIC NODE
//   CONNECTIVITY
//  ==============
//  Prnt is for a single parent (iff it fits conceptually) 
//         (typ."Up", but rendered to left, like in NexSTEP)
    prnt   *Nork  // L: One max, and at top of L-H list 
//  Kids is ordered child nodes and outgoing links 
//         (typ."Down", but render to right, like in NexSTEP)
    kids []*Nork  // R: Many, in order 
//  Usrs is referrers (like in ToC/ditamap) and incoming links
//         (typ."Up", but rendered to left, like in NexSTEP)
    usrs []*Nork  // L: Symlinks/Referrers, listed L-H, under Prnt (if has) 
//  Tags is metadata tags/facets (rendered above) (can include properties ?)
    tags []*Nork  // Tags/Facets
//  Vers is previous versions (is a list, not a tree)
//         (rendered below, as a stack of cards)
    vers []*Nork  // Versions
//  =================
//    ADVANCED NODE
//   CHARACTERISTICS
//    (FS-ORIENTED)
//  =================
    // isDir might be undefined for non-FS use cases; it can be more
    // generally defined as "canKids", i.e. is able to have kid nodes.
    isDir bool
    // isDirlike also includes symlinks (and other edge cases ?)
    // but should be considered TBD for anything but filesystem. 
    isDirlike bool
    // level starts at 0 for root, and isRoot() is (level == 0)
    // (For isRoot() we don't also/alternatively test on whether 
    // Prnt is nil, because we might find other uses for Prnt,
    // such as the file that contains a document tree.)
    // Discussion: It is equal to the number of "/" path separators
    // *separating* path elements (i.e. not including any leading or
    // trailing separators). Therefore it is 0 for an XML document
    // root node or the local root of a file & dir tree (where in
    // both cases, isRoot() is true and parent() is [probably] nil)), 
    // and it is >0 for others. Reserve negative numbers for future
    // (ab)use.
    // So a Root has a relFP of "." and an absFP that is the rooted 
    // absolute path of this root node w.r.t. the external environ-
    // ment (for a file or dir, the file system root; for a markup
    // node, the absolute path of the containing file.
    level int
    // ======================================
    //  Substructure for: Materialized Paths
    // ======================================
    // relPath is use case -specific, but often relative filepath.
    // Discussion: it is the relative path of this Nork, relative 
    // to its tree's root Nork, which is the "local root" shared 
    // with other Norks in the same interconnected tree. (That is 
    // to say, a local root is the highest/topmost node of a direc-
    // tory tree imported in a single batch.) The last element of 
    // the relFP is this Nork's own name/label, analagous to
    // FP.Base(Path).
    relPath string 
    // absPath is use case -specific, but often absolute filepath.
    // Discussion: It is the same as path, except that it is rooted 
    // in - i.e. it is traced back to the root of - a local file
    // system (or documwnt). For a file or dir in a filesystem,
    // it is rooted at the filesystem root. For a markup node
    // or a map/ToC file, it is rooted at the document start.
    absPath string 
//  =================
//    NODE IDENTITY,
//   USAGE, USE CASE 
//  =================
    Uuid7 string // So can handle ptrs btwn separate trees & lists 
    UC    string // QuadraMode // UC selector 
}

func (p *Nork) IsDir() bool  { return p.isDir }
func (p *Nork) Level()  int  { return p.level }
func (p *Nork) IsRoot() bool { return p.level == 0 }
func (p *Nork) IsDirlike() bool { return p.isDirlike }
func (p *Nork) AbsPath() string { return p.absPath }
func (p *Nork) RelPath() string { return p.relPath }
func (p *Nork) SetAbsPath(s string) { p.absPath = s }
func (p *Nork) SetRelPath(s string) { p.relPath = s }


