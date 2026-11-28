package nork

import(
    "github.com/fbaube/oid"
	FU "github.com/fbaube/fileutils"
	TL "github.com/fbaube/treelinks"
)

// StringFunc is used by interface Norker, so a 
// method signature actually (MAYBE!) looks like:
//   func (*Nork) FuncName() string
type StringFunc func(Norker) string

// LinksPtr overrides the same-named func of
// the *Treelinks, aktho it does the same thing. 
func (p *Nork) LinksPtr() *TL.Treelinks {
      return p.links
}

// Nork is an all-weather, all-terrain tree node.
// It adds kid node prdering on top of [Treelinks],
// and it provides for multiple sets of Treelinks
// for support advanced "Quadranode" funtionality
// like incoming references, attached metadata
// facets, and versioning. 
//
// Child nodes ("Kid"s) have an externally-defined 
// order, implying that child nodes should not be
// directly compared (via [Comparable]) amongst
// themselves.
//
// The ability to impose a specific specified order
// on child nodes is essential for representation
// of content markup.
//
// In addition to a basic set of Treelinks, several
// additional fields are provided, and documented
// here to make them visible, because many of these
// fields are not public (i.e. capitalised). 
//  - Name string: the typical field for a tree node 
//  - links Treelinks: links built & kept in-memory 
//  - oidLinks Treelinks: links as persisted
//  - OID oid.OID: includes the DB table, so that 
//    [OuterOwner] (the embedding super-struct) can 
//    be fetched by OID directly from the DB, rather 
//    than going thru some kind of in-memory router / 
//    map / OID registry that figures out which table.
//  - Errer: collection point for errors associated
//    with this tree node 
//  - UC string: Use Case, when using as a QuadraMode
//  - creatPath string: bookkeeping
//  - IsFocused, IsSelected, IsExpanded, IsVisible,
//    IsMafked bool: GUI-related stuff 
//  - level int: this is actually stored with Treelinks,
//    and is not appropriate here; note that zero marks
//    a root node; note that we don't also/alternatively
//    test on whether Prnt is nil, because we might find
//    other uses for Prnt, such as the file that contains
//    a document tree.) 
//    Discussion of level: It is equal to the number of
//    "/" path separators *separating* path elements (i.e.
//    not including any leading or trailing separators).
//    Therefore it is 0 for an XML document root node or
//    the local root of a file & dir tree (where in both
//    cases, isRoot() is true and parent() is [probably]
//    nil)), and it is >0 for others. Reserve negative
//    numbers for future (ab)use.
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
// UC is Use Case is QuadraMode UC selector
   UC    string 
// OID includes the DB table, so that [OuterOwner]
// (the embedding super-struct) can be fetched by 
// OID directly from the DB, rather than going 
// thru some kind of in-memory router / map / 
// OID registry that figures out which table.
   oid.OID

// ==============
//   BASIC NODE
//  CONNECTIVITY
// ==============
// ON-DISK and IN-MEM have the same structure, so maybe
// they could be redefined using generics. Each is a
// slice except Prnt Prev Next, which are a single items.
//  - Prnt is for a single parent (iff it fits conceptually) 
//    (typ."Up", but gets rendered to the left, like in NexSTEP)
//  - Kids is ordered child nodes and outgoing links 
//    (typ."Down", but gets rendered to the right, like in NexSTEP)
//  - Usrs is referrers (like in a ToC or ditamap) and incoming links
//    (typ."Up", but gets rendered to the left, like in NexSTEP)
//  - Tags is metadata tags/facets (gets rendered above)
//    (can include properties ?)
//  - Vers is previous versions (gets rendered below, as a deck)
//   (done as a list, not a tree)

// ===============
//  UUIDS ON-DISK
// ===============
// These OID are behind a pointer because if a
// tree is IN-MEM and its nodes do not correspond
// to objects ON-DISK, such as nodes in a XML file, 
// then assigning UUIDs is just wasteful and confusing.
   oidLinks *TL.Treelinks 
// =============
//  OIDS IN-MEM
// =============
   links *TL.Treelinks

// =============
//  OTHER ITEMS 
// =============

// Errer is here because a Nork that was created
// successfully and linked into a tree might later 
// acquire an error due to an external change.
FU.Errer 

// level probably does not belong here, but let's
// put it here to avoid a lot of rewrite (2026.05).
// level starts at 0 for root, and isRoot() iff
// (level == 0)
//
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
// A Root has a RelPath of "." and can (in principle) be an absFP 
// that is the rooted absolute path of this root node w.r.t. the
// external environment (for a file or dir, the file system root;
// for a markup node, the absolute path of the containing file).
   level int
   
// -----------
//  GUI Stuff
// -----------
   IsFocused, IsSelected, IsExpanded, IsVisible, IsMarked bool
}

// NEED ECHO (assume html), INFOS, DEBUG
/*
func (p *Nork) IsDir() bool  { return p.IsDir() }
func (p *Nork) IsRoot() bool { return p.level == 0 }
func (p *Nork) IsDirlike() bool { return p.IsDirlike() }
*/

func (p *Nork) CreatPath() string { return p.creatPath } 

// A Root has a RelPath of "." and can (in principle) be an absFP 
// that is the rooted absolute path of this root node w.r.t. the
// external environment (for a file or dir, the file system root;
// for a markup node, the absolute path of the containing file).
