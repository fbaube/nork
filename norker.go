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
	// node, it is meaningless except for a [RootNorker],
	// where it is the rel.path to the containing document.
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
	Root() TL.Treelinkser 
	IsRoot() bool
	IsDir() bool
	IsDirlike() bool
	// ============
	//  MISCELLANY
	// ============
	LinePrefixString() string
	PrintTree(io.Writer) error
}


