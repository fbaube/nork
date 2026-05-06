package nork

import( 
	L  "github.com/fbaube/mlog"
	FU "github.com/fbaube/fileutils"
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

/*

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

*/

func (p *NorkFactory) RootPath() string {
     return p.RootFPs.AbsFP
     }

// NewFilepathNork expects a relative path (but why? - 
// the reason is forgotten - but it does get thru more
// runtime security checks that way). It does not load 
// file content, because it is an expensive operation 
// that can and should be done elsewhere.
//
// It really only loads the ´Filepaths`, so it is kind
// of useless. It does not touch other fields in the `Nork`.
//
// Per the func name, do not use this for a path that
// is not a filesystem path.
// .
func (pFac *NorkFactory) NewFilepathNork(aRelPath string) *Nork {
	if aRelPath == "" {
		L.L.Error("NewFPnork: missing path")
		return nil 
	}
	// Note that this also allocates the Nork, and 
	// "should" provide access to its non-public fields 
	pNN := new(Nork)
	var e error 
	pNN.FPs = *FU.NewFilepaths(aRelPath)
	if pNN.FPs.HasError() {
	     	L.L.Error("NewFPnork: " + e.Error())
                return nil
        }
	return pNN
}

