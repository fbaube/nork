package nork

import(
	"fmt"
	L  "github.com/fbaube/mlog"
	FU "github.com/fbaube/fileutils"
)

// NewFSOTreeNork expects a relative path (but why? - the reason is
// forgotten - but it does get thru more runtime security checks 
// that way). It does not load file content, because it is an
// expensive operation that can and should be done elsewhere.
//
// It loads the ´Filepaths` which does perform a number of checks
// specific to filesystem objects. Then the content can be loaded
// lazily at any time.
// 
// Per the func name, do not use this for a path that is not
// a filesystem path.
// .
func (pFac *FSOTreeNorkFactory) NewFSOTreeNork(aRelPath string) *FSONork {
	if aRelPath == "" {
		L.L.Error("newfsotreenork: missing path")
		return nil 
	}
	// var e error 
	var pFN  *FSONork
	pFN = new(FSONork)
	pFN.Nork = *NewNork(aRelPath)
	pFN.FSO  = FU.NewFSObject(aRelPath)
	if pFN.Nork.HasError() {
	     	pFN.SetError(fmt.Errorf("newfsontreeork:newnork: %w",
			pFN.Nork.GetError()))
                return pFN
        }
	if pFN.FSO.HasError() {
	     	pFN.SetError(fmt.Errorf("newfsotreenork:newfso: %w",
                        pFN.FSO.GetError()))
                return pFN
        }
	// Now the [Filepaths] struct [FSO.FPs] is valid, and it has
	// the object's [os.FileInfo], and it can implement the full
	// range if functions that are based on the value of the
	// `FileInfo`, and we can ignore the field [Nork.CreatPath].

	// So here we could do things that depend upon and/or affect 
	// the internal state of the factory.

	return pFN
}

// NewFSOLoneNork expects a relative path (but why? - the reason is
// forgotten - but it does get thru more runtime security checks 
// that way). It does not load file content, because it is an
// expensive operation that can and should be done elsewhere.
//
// It loads the ´Filepaths` which does perform a number of checks
// specific to filesystem objects. Then the content can be loaded
// lazily at any time.
// 
// Per the func name, do not use this for a path that is not
// a filesystem path.
// .
func NewFSOLoneNork(aRelPath string) *FSONork {
	if aRelPath == "" {
		L.L.Error("newfsolonenork: missing path")
		return nil 
	}
	// var e error 
	var pFN  *FSONork
	pFN = new(FSONork)
	pFN.Nork = *NewNork(aRelPath)
	pFN.FSO  = FU.NewFSObject(aRelPath)
	if pFN.Nork.HasError() {
	     	pFN.SetError(fmt.Errorf("newfsonloneork:newnork: %w",
			pFN.Nork.GetError()))
                return pFN
        }
	// IsDir ??
	if pFN.FSO.HasError() {
	     	pFN.SetError(fmt.Errorf("newfsolonenork:newfso: %w",
                        pFN.FSO.GetError()))
                return pFN
        }
	// Now the [Filepaths] struct [FSO.FPs] is valid, and it has
	// the object's [os.FileInfo], and it can implement the full
	// range if functions that are based on the value of the
	// `FileInfo`, and we can ignore the field [Nork.CreatPath].

	// So here we could do things that depend upon and/or affect 
	// the internal state of the factory.

	return pFN
}

