package nork

import(
	FU "github.com/fbaube/fileutils"
)

// FSONork embeds a Nork "QuadraNode".
//
// The UC (Use Case) is:
//  - UC.FSI: File System Item: represent a dir-or-file-or-softlink:
//    	      Here ordering is less important. Paths are valid FS paths. 
//            Prnt:dir, Kids:contents(!dir-v-file), Usrs:incoming-symlinks 
//
// Using Norks for files & dirs exhibits strong typing. Dirs are dirs
// and files are files and never the twain shall meet. This means that
// (a) dirs cannot contain own-content, and (b) files can never be
// non-leaf nodes. (Note tho that symlinks have aspects of both.)
// However this dir/file/etc typing is too complex to handle here 
// in a Nork, because a leaf node can be either a file or a dir, 
// and a field like "canKid bool" is a bit OTT, so the file/dir 
// distinction is handled instead by an outer struct type that 
// embeds Nork, such as FSONork, embedding [fileutils.FSObject].
//
// If we build up a tree of Norks when processing an [os.DirFS], the
// strict ordering provided by DirFS is not strictly needed, BUT it
// can anyways be used (and relied upon) because the three flavors of
// WalkDir are deterministic (using lexical order). WalkDir does tho
// promise that a given Nork will always appear AFTER the Nork for
// its directory has appeared, which makes it "easy" to build a tree.
// .
type FSONork struct {
    Nork
    // FSO includes paths and flags like IsDir and DoesNotExist.
    // (For non-FS use case, IsDir() might be sorta "CanKids".)
    // (IsDirlike should be considered TBD for anything but FS.)
    FSO *FU.FSObject
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
}

// func (p *FSONork) IsDir()  bool { return p.Nork.IsDir() }
func (p *FSONork)   IsRoot()  bool { return p.Nork.Level() == 0 }
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

