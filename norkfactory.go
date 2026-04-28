package nork

import(
	"errors"
	L  "github.com/fbaube/mlog"
	FU "github.com/fbaube/fileutils"
	FP "path/filepath"
)
	
// NorkFactory creates a root node and then stores and tracks the 
// state of a Nork tree being assembled via calls to its method 
// [NewNork]. This seems like a good idea because then we can spin 
// up Nork factories that are customised for any of many purposes. 
type NorkFactory struct {
	// nexSeqID should be unique across all Norks, 
	// so it should probably be filled in by SQLite.
	nexSeqID      int
	RootFPs	      *FU.Filepaths 
//	isDir	      bool 
//	summaryString StringFunc
}

// NewNorker will make it possible to supply 
// a custom `New` func to a new factory. 
// type NewNorker func(string) (Norker, error)

// NewNorkFactory assumes (for now) that it is working with
// a filesystem; it verifies that it got a directory, sets 
// the bool [isDir], and creates and returns the root Nork.
// TBD: Global index counter or ask SQLite ? 
func NewNorkFactory(rootPath string) (*NorkFactory, *Nork, error) {
     	var s string 
        var e error
	var pNF *NorkFactory
	// Check the path 
	if rootPath == "" {
	   	s = "NewNorkFactory: missing root path"   	    
		L.L.Error(s) 
		return nil, nil, errors.New(s)
	}
	pNF = new(NorkFactory)
	pNF.RootFPs, e = FU.NewFilepaths(rootPath) 
	if e != nil || !pNF.RootFPs.IsDir { 
	   	s = "NewNorkFactory: bad or not-a dir: " + rootPath 
		L.L.Error(s)
		return nil, nil, errors.New(s)
	}
	// FIXME? pNF has to be valid here !!
	pRootNork := pNF.NewNork(rootPath)
	if pRootNork == nil {
	   s = "NewNorkFactory: cannot make root nork: " + rootPath
	   L.L.Error(s)
	   return nil, nil, errors.New(s)
	}
	L.L.Debug("RootNork's absFP: " + pRootNork.AbsPath())
	L.L.Debug("RootNork's relFP: " + pRootNork.RelPath())

	// For the relative path, try to trim the entire
	// RootNork RootPath off of this absolute path.
	// func CutPrefix(s, prefix string) (after string, found bool):
	// It returns s without the provided leading prefix 
	// string and reports whether it found the prefix.
	// If s dusn't start with prefix, CutPrefix returns (s, false).
	// If prefix is the empty string, CutPrefix returns (s, true).
	
	pRootNork.SetRelPath(pNF.RootFPs.AbsFP)
	// p.isRoot = true // zero value of Level is OK 
	pNF.RootFPs.IsDir = true // be sure 
	return pNF, pRootNork, nil
}

func (p *NorkFactory) RootPath() string {
     return p.RootFPs.AbsFP
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

