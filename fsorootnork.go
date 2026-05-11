package nork

import(
	FU "github.com/fbaube/fileutils"
)

type FSORootNork struct {
     	// nexSeqID should be unique across all Norks, 
	// so it should probably be filled in by SQLite.
        nexSeqID      int
	RootNork      *FSONork
	BatchID       int 
	FU.FSObject
}

