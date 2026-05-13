// nork is a Node with Ordered Kids. It stores basic
// bidirectional parent/kid relationships, plus other 
// useful hierarchy-related info like levels and paths. 
//
// interface [Norker] is implemented not for struct
// [Nork] but rather for the pointer, i.e. `*Nork´. 
// This makes nodes writable and also sharable.
//
package nork

import(
	FU "github.com/fbaube/fileutils"
	// "uuid"
	// func NewV7() UUID
)

// StringFunc is used by interface Norker, so a 
// method signature actually (MAYBE!) looks like:
//   func (*Nork) FuncName() string
type StringFunc func(Norker) string

// Nork is a "QuadraNode", handling four directions of
// connections for four distinct node user cases (UCs). 
// Nork is kind of a maximalist implementation, and is
// subject to redefinition and getting slimmed down.
//
// Paths need not be in the filesystem, and so paths are
// IAW package [path] rather than package [path/filepath].
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
// Note that UUID might take on outsized importance, because it might be
// used to reference all manner of other DB tables and data structures
// and external (to the system) targets. 
// It may become necessary to add prefixes to UUIDs to specify (e.g.) 
// the DB tables of targets.
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
// ===============
//  NODE IDENTITY
//   AND USE CASE 
// ===============
// creatPath is the path used to instantiate the node.
// It is not validated, and should conform to package 
// [path] (not [path/filepath]), and should be ignored 
// after another path is available in an embedding struct.
// It always uses forward slash "/" as the path separator.
// Note that we won't also define an AbsPath here in Nork,
// because then its presence thruout a tree makes moving 
// the root of the tree (i.e.. changing its absolute path) 
// too damned work-intensive and error-prone.
   creatPath string 
// Name is the only field that we know works across
// all Use Cases, and any rules are specific to UC. 
   Name  string
// Uuid7 is so that we can make & follow links btwn separate trees & lists.
   Uuid7 string 
// UC is Use Case is QuadraMode UC selector
   UC    string 
// ==============
//   BASIC NODE
//  CONNECTIVITY
// ==============
// Prnt is for a single parent (iff it fits conceptually) 
//        (typ."Up", but rendered to left, like in NexSTEP)
   prnt   *Nork  // L: One max, and at top of L-H list 
// Kids is ordered child nodes and outgoing links 
//        (typ."Down", but render to right, like in NexSTEP)
   kids []*Nork  // R: Many, in order 
// Usrs is referrers (like in ToC/ditamap) and incoming links
//        (typ."Up", but rendered to left, like in NexSTEP)
   usrs []*Nork  // L: Symlinks/Referrers, listed L-H, under Prnt (if has) 
// Tags is metadata tags/facets (rendered above) (can include properties ?)
   tags []*Nork  // Tags/Facets
// Vers is previous versions (is a list, not a tree)
//        (rendered below, as a stack of cards)
   vers []*Nork  // Versions
// ==========
//  And also
// ==========
// Errer is here because a Nork that is already linked into 
// a tree might acquire an error due to an external change.
FU.Errer 
// level probably does not belong here, but let's
// put it here to avoid a lot of rewrite (2026.05).
// 
// level starts at 0 for root, and isRoot() <=> (level == 0)
// (For isRoot() we don't also/alternatively test on whether
// Prnt is nil, because we might find other uses for Prnt,
// such as the file that contains a document tree.)
// 
// Discussion: It is equal to the number of "/" path separators
// - *separating* path elements (i.e. not including any leading 
// or trailing separators). Therefore it is 0 for an XML document
// root node or the local root of a file & dir tree (where in 
// both cases, isRoot() is true and parent() is [probably] nil)),
// and it is >0 for others. Reserve negative numbers for future
// (ab)use.
// 
// A Root has a RelPath of "." and can (in principle) an absFP 
// that is the rooted absolute path of this root node w.r.t. the
// external environment (for a file or dir, the file system root;
// for a markup node, the absolute path of the containing file).
   level int
// -----------
//  GUI Stuff
// -----------
   HasFocus, IsSelected, IsExpanded, IsVisible bool
}

// NEED ECHO (assume html), INFOS, DEBUG
/*
func (p *Nork) IsDir() bool  { return p.IsDir() }
func (p *Nork) IsRoot() bool { return p.level == 0 }
func (p *Nork) IsDirlike() bool { return p.IsDirlike() }
*/

func (p *Nork) Level()  int  { return p.level }
func (p *Nork) SetLevel(i int) { p.level = i }
func (p *Nork) CreatPath() string { return p.creatPath } 

// ===========================
//  Abs.Path and Rel.Path are
//   for: Materialized Paths
// ===========================
/*
// AbsPath is use case -specific, but often an absolute filepath.
// Discussion: It is the same as rel.path, but expanded to be
// rooted in - i.e. traced back to the root of - a local file
// system (or documwnt). For a file or dir in a filesystem,
// it is rooted at the filesystem root. For a markup node
// or a map/ToC file, it is rooted at the document start.
func (p *Nork) AbsPath() string { return p.AbsPath() }
func (p *Nork) SetAbsPath(s string) { p.SetAbsPath(s) }

// RelPath is use case -specific, but often a relative filepath.
// Discussion: it is the relative path of this Nork, relative to
// its tree's root Nork, which is the "local root" shared with
// other Norks in the same interconnected tree. (That is to say,
// a local root is the highest/topmost node of a directory tree
// imported in a single batch.) The last element of the relFP 
// is this Nork's own name/label, analagous to FP.Base(Path).
func (p *Nork) RelPath() string { return p.RelPath() }
func (p *Nork) SetRelPath(s string) { p.SetRelPath(s) }
*/
