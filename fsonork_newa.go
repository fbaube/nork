package nork

import(
	"os"
	"errors"
	FU "github.com/fbaube/fileutils"
)

// Arg s is a filepath. [FSONork] embeds [Errer].
// If (Errer.HasError], it is a [*os.PathError].
func newFSOnork(s string) *FSONork {
        var pFN = new(FSONork)
	if s == "" {
		pFN.SetError(errors.New("newfsonork: missing path"))
		return pFN
	}
        pFN.Nork = *NewNork(s)
        pFN.FSO = FU.NewFSObject(s)
	var pPE = new(os.PathError{Path:s})
        if pFN.Nork.HasError() {
	   pPE.Op = "newfsonork:newnork"
	   pPE.Err = pFN.Nork.GetError()
	   pFN.SetError(pPE)
        }
        if pFN.FSO.HasError() {
	   pPE.Op = "newfsonork:newfsonork"
	   pPE.Err = pFN.FSO.GetError()
	   pFN.SetError(pPE)
        }
	return pFN
 }

// NewFSOTreeNork prefers a relative path. It does not load file 
// content, which is an expensive operation done elsewhere. But it
// loads the [Filepaths], which does checks specific to filesystem
// objects. Then the content can be loaded lazily at any time.
// 
// Do NOT use this for a path that is not a filesystem path.
// .
func (pFac *FSOTreeNorkFactory) NewFSOTreeNork(aRelPath string) *FSONork {
	var pFN = newFSOnork(aRelPath)
	// If error, quick return. 
	if  pFN.FSO.HasError() {
                return pFN
        }
	// Now the [Filepaths] struct [FSO.FPs] is valid, and it has
	// the object's [os.FileInfo], and it can implement the full
	// range of functions that are based on the value of the
	// [os.FileInfo], and we can ignore the field [Nork.CreatPath].

	// So here we could do things that depend upon and/or affect 
	// the internal state of the factory.

	return pFN
}

// NewFSOLoneNork prefers a relative path. It does not load file
// content, which is an expensive operation done elsewhere. But it
// loads the [Filepaths], which does checks specific to filesystem
// objects. Then the content can be loaded lazily at any time.
//
// Do NOT use this for a path that is not a filesystem path.
// .
func NewFSOLoneNork(aRelPath string) *FSONork {
	var pFN = newFSOnork(aRelPath)
	// If error, quick return. 
	if  pFN.FSO.HasError() {
                return pFN
        }
	// Now the [Filepaths] struct [FSO.FPs] is valid, and it has
	// the object's [os.FileInfo], and it can implement the full
	// range of functions that are based on the value of the
	// [os.FileInfo], and we can ignore the field [Nork.CreatPath].

	// IsDir ??
	return pFN
}

