package nork

import(
	FU "github.com/fbaube/fileutils"
	TL "github.com/fbaube/treelinks"
)

// StringFunc is used by interface Norker, so a 
// method signature actually (MAYBE!) looks like:
//   func (*Nork) FuncName() string
type StringFunc func(Norker) string

// Nork is an all-weather, all-terrain tree node. 
//
// Child nodes ("Kid"s) have an externally-defined 
// order, implying that child nodes are not directly 
// [Comparable] amongst themselves. This specific 
// specified order of child nodes is essential
// for representation of content markup.
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
// Uuid7 is so that we can make & follow links btwn
// separate trees & lists. A Uuid7 includes a DB
// table name and (probably?) a DB rowID. 
   Uuid7
// UC is Use Case is QuadraMode UC selector
   UC    string 

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
/*
// ================
//  ON-DISK: UUIDs
// ================
   // These "u" (UUID) are behind a pointer because if
   // a tree is IN-MEM and its nodes do not correspond
   // to objects ON-DISK, such as nodes in a XML file, 
   // then assigning UUIDs is just wasteful and confusing.
   uLinks *Treelinks 
   /*
   prntU   Uuid7  // Left: One max, and at top of L-H list 
   kidsU []Uuid7  // Rite: Many, in order 
   usrsU []Uuid7  // Left: Symlinks/Referrers
   tagsU []Uuid7  // Topp: Tags/Facets
   versU []Uuid7  // Botm: Versions
   */
// ===============
//  IN-MEM: *Nork 
// ===============
   // These "m" (memory) proabbly do not need to be a 
   // behind pointer, but uLinks does it so let's here too. 
   *TL.Treelinks
   /*
   prnt   *Nork  // Left: One max, and at top of L-H list 
   kids []*Nork  // Rite: Many, in order 
   usrs []*Nork  // Left: Symlinks/Referrers
   tags []*Nork  // Topp: Tags/Facets
   vers []*Nork  // Botm: Versions
   */
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
   IsFocused, IsSelected, IsExpanded, IsVisible bool
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
