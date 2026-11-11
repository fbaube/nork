package nork

import (
	"io"
    TL "github.com/fbaube/treelinks"
)

type Norker interface {
     	TL.Stringser
        TL.Treelinkser 
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
	// SetLevel(int) // should be self-calculated by node 
	// Root should always return the root
	//  (at arena index 0, if one is used) 
	Root() TL.Treelinkser // RootNork
	IsRoot() bool
	IsDir() bool
	IsDirlike() bool
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
	SetFirstKid(TL.Treelinkser)
	SetLastKid (TL.Treelinkser)
	SetPrevPeer(TL.Treelinkser)
	SetNextPeer(TL.Treelinkser)
}



