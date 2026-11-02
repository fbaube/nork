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

// Nork is TBS.
//
// Child nodes ("Kid"s) have an externally-defined order,
// implying that child nodes are not directly [Comparable] 
// amongst themselves. This specific specified order or 
// child nodes is essential for representation of content.
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
