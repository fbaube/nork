package nork

import(
	FU "github.com/fbaube/fileutils"
)

// FSONork embeds an FSO with a Nork (node with ordered kids). 
// The FSO includes a [Filepaths] struct and flags like IsDir
// and DoesNotExist.
//  - For a non-FS use case, IsDir() might be a "CanKids"
//  - For any non-FS use case, IsDirlike should be considered TBD 
//
// The UC (Use Case) for FSONork is:
//  - UC.FSI: File System Item: represent a dir-or-file-or-softlink:
//    Here ordering is less important. Paths are valid FS paths. 
//    prnt-link is a dir, kid links are dir-or-file-or-symlink,
//    w content; usrs-linnks (if used) can be incoming-symlinks. 
//
// Using Norks for files & dirs is asking too much of them. For
// one thing, symlnks add a lot of complexity. This complexity 
// is instead handler by this struct, FSONork, an outer struct
// type that embeds a Nork and an [FSObject].
//
// Implementation node: If we build up a tree of Norks when
// processing an [os.DirFS], the strict ordering provided by
// DirFS is not strictly needed, BUT it can anyways be used
// (and relied upon) because the three flavors of WalkDir 
// are deterministic (using lexical order). WalkDir does 
// tho at least promise that a given Nork will always appear 
// only after the Nork for its directory has appeared, which
// makes it "easy" to build a tree.
// .
type FSONork struct {
    Nork
    FSO *FU.FSObject
}

// Why is this commented out ? Name collision ? 
// func (p *FSONork) IsDir()  bool { return p.Nork.IsDir() }
// This should check level instead. 
func (p *FSONork)   IsRoot()  bool { return p.Nork.links.Parent() == nil }
func (p *FSONork) IsDirlike() bool { return p.FSO.IsDirlike() }

// ===========================
//  Abs.Path and Rel.Path are
//   for: Materialized Paths
// ===========================

// AbsFP is an absolute filepath.
// Discussion: It is the same as rel.path, but expanded to be
// rooted in - i.e. traced back to the root of - a filesystem.
func (p *FSONork) AbsFP() string { return p.FSO.FPs.AbsFP }
// func (p *FSONork) SetAbsFP(s string) { p.FSO.FPs.AbsFP = s }

// RelPath is a relative filepath.
// Discussion: it is the relative path of this Nork, relative to
// its tree's root Nork, which is the "local root" shared with
// other Norks in the same interconnected tree. (That is to say,
// a local root is the highest/topmost node of a directory tree
// imported in a single batch.) The last element of the relFP
// is this Nork's own name/label, analagous to FP.Base(Path).
func (p *FSONork) RelFP() string { return p.FSO.FPs.RelFP }
// func (p *FSONork) SetRelFP(s string) { p.FSO.FPs.SetRelFP = s }

// ==================

// A Root has a RelPath of "." and can (in principle) be an absFP 
// that is the rooted absolute path of this root node w.r.t. the
// external environment (for a file or dir, the file system root;
// for a markup node, the absolute path of the containing file).
// level int
   
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

