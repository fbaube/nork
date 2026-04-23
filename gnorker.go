package nork

import (
	"io"
)

// Stringser is copied from
// [github.com/fbaube/stringutils.Stringser] to keep out dependencies.
type Stringser interface {
        Echo()  string
        Infos() string
        Debug() string
}

// The interfaces have to be divided into two sets:
//  - One for Nork, that works with a basic tree.
//  - One for  Cnty, that works for such a tree, but
//    with additional func's defined for a filesystem.
//  - Cnty can redefine-and-override Nork versions. 

// Norker is satisfied by [*Nork] NOT by Nork.
type Norker interface {
     	Stringser
	// =====================
	//  PATHS (Rel. & Abs.) 
	// =====================
	// RelPath is rel.filepath for a file/dir, and for a DOM
	// node, it is meaningless, unless it is a [RootNorker],
	// for which it is the rel.path to the containing document. 
	RelPath() string
	SetRelPath(string) 
	// AbsPath is abs.filepath for a file/dir, and for a DOM
	// node, the (abs.)path of the node w.r.t. the document
	// root, except for a [RootNorker], for which it is
	// the abs.path to the containing document. 
	AbsPath() string
	SetAbsPath(string)
	// ======================
	//  PATH CHARACTERISTICS
	// ======================
	// Level is zero-based (i.e. root Nork's is 0) 
	Level() int
	// PACKAGE METHODS
	// SetLevel(int) // should be self-calculated by node 
	// Root should always return the root
	//  (at arena index 0, if one is used) 
	Root() *Nork // RootNork
	IsRoot() bool
	IsDir() bool
	IsDirlike() bool
	// ==================
	//   NODE LINKS i.e.
	//  INTERCONNECTIONS
	// ==================
	Parent() *Nork
	SetParent(*Nork)
	HasKids() bool
	KidsAsSlice() []*Nork
	FirstKid() *Nork
	LastKid()  *Nork
	PrevPeer() *Nork // cos "PrevKid" = iterator on the parent 
	NextPeer() *Nork
	// AddKid returns the just-added kid, who
	// is added at the end of the list of Kids,
	// and who (TODO FWIW) knows his own arena
	// index (using [slices.Index])
	AddKid(*Nork) *Nork
	// AddKids returns the method target
	// - the parent of all the kids 
	AddKids([]*Nork) *Nork
	ReplaceBy(*Nork) *Nork // returns Whom
	// ============
	//  MISCELLANY
	// ============
	LinePrefixString() string
	PrintTree(io.Writer) error
}

// Cntyer is satisfied by [*Cnty] NOT by Cnty.
type Cntyer interface {
	// =========================
	//  FILEPATHS (Rel. & Abs.) 
	// =========================
	// RelFP is rel.filepath for a file/dir.
	RelFP() string
	SetRelFP(string) 
	// AbsFP is abs.filepath for a file/dir.
	AbsFP() string
	SetAbsFP(string)
	// Kids stuff that only makes sense for Kids in a linked list. 
	SetFirstKid(*Nork)
	SetLastKid (*Nork)
	SetPrevPeer(*Nork)
	SetNextPeer(*Nork)
}

