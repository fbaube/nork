package nork

import(
	"os"
	"errors"
	L  "github.com/fbaube/mlog"
)

// FSOTreeNorkFactory stores and tracks the state of an FSO 
// (File System Object) Nork tree being assembled via calls
// to its method [NewFSOTreeNork].
//
// It cannot (and should not) become invalid after creation,
// so Errer is not embedded & used.
// 
// A Factory is a good idea because then we can add arbitrary
// additional functionality (and complexity) and also spin up
// Nork factories that are customised for any of many purposes. 
// .
type FSOTreeNorkFactory struct {
	// NexSeqID should be unique across all Norks, 
	// so it should probably be filled in by SQLite.
	NexSeqID      int
	// RootNork includes the [Filepaths].
	RootFSONork  *FSONork
	// BatchID indexes a DB table. 
	BatchID	      int 
}

// NewNorker will make it possible to supply 
// a custom `New` func to a new factory. 
// type NewNorker func(string) (Norker, error)


// NewFSOTreeNorkFactory works with a directory in a filesystem:
//  - it verifies that it got a directory 
//  - it also creates and returns the root Nork
//  - an error return is a [*os.PathError]
// 
// If the target is a single file (instead of a directory),
// just use func [NewFSOLoneNork], which is also used to
// create the initial root Nork in this func.
//
// TBD: Global index counter or ask SQLite ?
// TBD: Maybe also pass in a "NewNork" function.
// .
func NewFSOTreeNorkFactory(aRootPath string) (
       *FSOTreeNorkFactory, *FSONork, error) {
     
     	var s string 
	var pFTNF = new(FSOTreeNorkFactory)
	var pPE  = new(os.PathError { Path:aRootPath })
	// Check the path 
	if aRootPath == "" {
	   	pPE.Op = "newfsotreenorkfac"
		pPE.Err = errors.New("missing root path")
		return pFTNF, nil, pPE
	}
	pFTNF.RootFSONork = NewFSOLoneNork(aRootPath)
	if pFTNF.RootFSONork.HasError() || !pFTNF.RootFSONork.FSO.FPs.IsDir {
	   	pPE.Op = "newfsotreenorkfac.newroot"
		pPE.Err = errors.New("not-a dir or bad dir")
		return pFTNF, nil, pPE
	}	
	// FIXME? pFTNF has to be valid here !!
	var pRootFSONork *FSONork
	pRootFSONork = NewFSOLoneNork(aRootPath)
	if pRootFSONork == nil {
	   s = "NewFSONorkFactory: cannot make root nork: " + aRootPath
	   L.L.Error(s)
	   return nil, nil, errors.New(s)
	}
	L.L.Debug("RootFSONork's absFP: " + pRootFSONork.AbsFP())
	L.L.Debug("RootFSONork's relFP: " + pRootFSONork.RelFP())

	// For the relative path, try to trim the entire
	// RootFSONork RootPath off of this absolute path.
	// func CutPrefix(s, prefix string) (after string, found bool):
	// It returns s without the provided leading prefix 
	// string and reports whether it found the prefix.
	// If s dusn't start with prefix, CutPrefix returns (s, false).
	// If prefix is the empty string, CutPrefix returns (s, true).
	
//	pRootFSONork.SetRelPath(pFTNF.RootFSONork.FPs.AbsFP)
//	pRootFSONork.RelFP = pFTNF.FSONork.FSO.AbsFP
	// p.isRoot = true // zero value of Level is OK 
//	pFTNF.RootFSONork.FPs.IsDir = true // be sure 
	return pFTNF, pRootFSONork, nil
}

func (p *FSOTreeNorkFactory) RootPath() string {
     return p.RootFSONork.FSO.FPs.AbsFP
     }



