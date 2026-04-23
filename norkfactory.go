package nork

import(
	"errors"
	L  "github.com/fbaube/mlog"
	FU "github.com/fbaube/fileutils"
	FP "path/filepath"
)
	
// NorkFactory creates a root node and then stores and tracks the state
// of a Nork tree being assembled via calls to its method [NewNork].
type NorkFactory struct {
	// nexSeqID should be unique across all Norks, 
	// so it should probably be filled in by SQLite.
	nexSeqID      int
	rootPath      string
	isDir	      bool 
//	summaryString StringFunc
}

// NorkEng is a package global, which is dodgy and not re-entrant.
// Axctually now a global integer for primary keys might work. 
// var NorkEng *NorkEngine = new(NorkEngine)

type NewNorker func(string) (Norker, error)

// NewNorkFactory verifies it got a directory, and then sets 
// the bool [isDir]. TBD: Global index counter or ask SQLite ? 
func NewNorkFactory(rootPath string) (*NorkFactory, error) {
     	var s string 
        // var e error
	var pGF *NorkFactory
	
	// L.L.Debug("NewNorkFactory: starting seqID: %d", NorkEng.nexSeqID)
	if rootPath == "" {
	   	s = "NewNorkFactory: missing root path"   	    
		L.L.Error(s) 
		return nil, errors.New(s)
	}
	pGF = new(NorkFactory)
	// NOTE the next stmts assume *filesystem* not XML DOM
	// FIXME This call does not guarantee it is an AbsFP
	// asAbsPath := FU.EnsureTrailingPathSep(FP.Clean(rootPath))

	// Verify that it is in fact a directory
	if !FU.IsDirAndExists(rootPath) { 
	   	s = "NewNorkFactory: path is not a dir: " + rootPath 
		L.L.Error(s)
		return nil, errors.New(s)
	}
	// FIXME pGF has to be valid here !!
	p := pGF.NewNork(rootPath)
	if p == nil {
	   s = "NewNorkFactory: cannot make node: " + rootPath
	   L.L.Error(s)
	   return nil, errors.New(s)
	}
	// CHECK THE PATHS
	L.L.Debug("RootNode's absFP: " + p.AbsPath())
	L.L.Debug("RootNode's relFP: " + p.RelPath())

	// For the relative path, try to trim the entire
	// RootNode RootPath off of this absolute path.
	// func CutPrefix(s, prefix string) (after string, found bool):
	// It returns s without the provided leading prefix 
	// string and reports whether it found the prefix.
	// If s dusn't start with prefix, CutPrefix returns (s, false).
	// If prefix is the empty string, CutPrefix returns (s, true).
	
	// p.relPath = p.absPath.S()
	p.SetRelPath(p.absPath)
	// p.isRoot = true // must use Level()==0
	// p.isDir = true  // FIXME 2026
	return pGF, nil
}

func (p *NorkFactory) RootPath() string {
     return p.rootPath
     }

// NewNork expects a relative path (!!), and does not either
// (a) set/unset the bool [isDir] or (b) load file content,
// because these are expensive operations that can and should
// be done elsewhere, and also (c) they do not apply if this
// is being used for XML DOM. 
func (pFac *NorkFactory) NewNork(aRelPath string) *Nork {
	if aRelPath == "" {
		L.L.Error("NewNork: missing path")
		return nil 
	}
	// Note that this also allocates the Nork, and 
	// "should" provide access to its non-public fields 
	pG := new(Nork) 
	pG.SetRelPath(aRelPath)
	asAbsPath := FP.Join(pFac.RootPath(), aRelPath)
	if FU.IsDirAndExists(asAbsPath) {
	   asAbsPath = FU.EnsureTrailingPathSep(asAbsPath)
	   }
	pG.absPath = asAbsPath // FU.AbsFilePath(asAbsPath) 
	// pG.isDir =... sorry, not done here 
	return pG
}

