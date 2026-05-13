package nork

import(
	"errors"
	"path"
	// "uuid" 
)

// NewNork does NOT assume an FS-type environment, and is fairly simple:
//  - It expects a relative path, not absolute; "." and "./" are OK but
//    not ""; an absolute path does NOT return an error 
//  - The path need not be in a filesystem, but paths are still expected 
//    to conform to package [path] (rather than package [path/filepath]):
//    "The path package should only be used for paths separated by forward
//    slashes, such as the paths in URLs. This package does not deal with
//    Windows paths with drive letters or backslashes; to manipulate
//    operating system paths, use the path/filepath package."
//  - It does not attempt to resolve the path to an absolute path:
//    this capability is not part of the `path` package 
//  - Its only concept is tree, not filesystem, so it mainly just records
//    its own node name and acquires a UUID7
//  - It implements interface [Errer], so errors can be acquired dynamically
//  - There is no distinction between a "tree" Nork and a "lone" Nork 
// .
func NewNork(aRelPath string) *Nork {
     	// Save away the creation-time path 
        var pN = new(Nork{creatPath:aRelPath, Name:path.Base(aRelPath)})
	if aRelPath == "" {
	   pN.SetError(errors.New("newnork: missing (relative) path"))
	}
	// FIXME pN.Uuid7 = uuid.NewV7() 
        return pN
}


