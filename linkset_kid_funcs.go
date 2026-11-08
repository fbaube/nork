package nork

import (
	"fmt"
	"os"
	"slices"
)

// All these func's use generics on << (p *Linkset[T]) >>, 
// and have to work when T is either [*Nork] or [Uuid7].

// Func's in the file are divided into three groups:
//  - Basic parent/child stuff
//     - Parent() T 
//     - SetParent(p2 T) 
//     - HasKids() bool 
//     - ValidateKids() err 
//  - Prev & Next peer (i.e. sibling node) 
//     - Prev() T 
//     - Next() T
//     - SetPrev(p2 T)
//     - SetNext(p2 T)
//  - Kid queries
//     - FirstKid() T
//     - LastKid() T
//     - KidsAsSlice() []T 
//  - Kid changes
//     - AddKids(rKids []T) T // returns p, the Parent
//     - AddKid (aKid   T) T // returns aKid
//     - SetFirstKid(p2 T)  // does what ?
//     - SetLastKid(p2 T)  // does what ?
//     - func's to delete kid(s) ? 

// 2026.11: We need a way to persist the order of Kid nodes 
// to the DB, so let's
//  - keep using fields Prev and Next, and
//  - make sure they get persisted to the DB, and
//  - when storing them in memory, continue to use a
//    linked list instead of storing kids in a slice.

// We'd like to use a slice for Kids. But it doesn't
// preserve prev/next info when persisted to DB, so
// when we read it back in we have to check paths 
// to match kids to parents, and for ordering, we
// have to rely on the order returned from the DB.
// This complexity is not robust!
//
// So it makes sense to match parents and kids (and
// prev and next) using UUIDs, and then the paths of
// parents and kids can be checked as a redundancy
// for finding logic errors.
//
// So when a tree is redd (back) in from the DB, 
// it should be straightforward to validate each
// set of (parent + kid(s)).
//
// 1-to-1 links can be stored dirctly in persisted
// items rather than in a línk table, as required
// for 1-to-N or N-to-1. Links for parent, prev,
// next, firstKid, and lastKid are all 1-to-1.
//
// Therefore a persisted Nork has five fields:
// Prnt, Prev, Next, Kid1, KidN.

// ==========================
//  BASIC PARENT/CHILD STUFF
// ==========================

// HasKids is duh.
func (p *Linkset[T]) HasKids() bool {
	// return p.FirstKid() != nil && p.LastKid() != nil
	return p.kids != nil && len(p.kids) != 0
}

// Parent returns the parent, duh.
// func (p *Linkset[T]) Parent() T {
func (p *Linkset[T]) Parent() T {
	return p.prnt
}

// SetParent has no side effects.
func (p *Linkset[T]) SetParent(p2 T) {
	p.prnt = p2
}

// SetPrev has no side effects.
func (p *Linkset[T]) SetPrev(p2 T) {
	p.SetPrev(p2)
}

// SetNext has no side effects.
func (p *Linkset[T]) SetNext(p2 T) {
	p.SetNext(p2)
}

/* N/A for SLICES (instead of linked lists)

// SetFirstKid has no side effects.
func (p *Linkset[T]) SetFirstKid(p2 T) {
	// p.SetFirstKid(p2)
	// https://stackoverflow.com/questions/53737435/how-to-prepend-int-to-slice
	p.kids = append([]*Nork{p2}, p.kids...)
}

// SetLastKid has no side effects.
func (p *Linkset[T]) SetLastKid(p2 T) {
	// p.SetLastKid(p2)
	p.kids = append(p.kids, p2) 
}

*/

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
// NOTE that just about ALL old (linked list) code is ripped out. 
func (p *Linkset[T]) AddKid(aKid *Nork) T { // returns aKid
	// fmt.Printf("nord: ptrs? aKid<%T> p<%T> \n", aKid, p)
	if aKid.Prev() != nil || aKid.Next() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("AddKid(K) can't cos K has siblings")
	}
	if aKid.Parent() != nil && aKid.Parent() != p {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("E.AddKid(K) can't cos K has non-P parent")
	}
	// Set the level & parent now
	aKid.SetLevel(p.Level() + 1)
	aKid.SetParent(p)
	if p.kids == nil || len(p.kids) == 0 { 
		// ---------------------------------------
		//  No Kids yet: New Kid will be ONLY Kid
		// ---------------------------------------
		p.kids = []*Nork{ aKid } 
	} else { 
		// --------------------------------
		//  So, APPEND as the new last kid
		// --------------------------------
		p.kids = append(p.kids, aKid)
	}
	return aKid
}

// AddKids adds the supplied nodes as kids, after 
// any pre-existing kids, and returns the parent.
// The supplied nodes may NOT be already linked
// to each other as a linked list. 
func (p *Linkset[T]) AddKids(rKids []*Nork) T { // returns p 
	// fmt.Printf("nord: ptrs? aKid<%T> p<%T> \n", aKid, p)
	for _, aKid := range rKids {	
	    if aKid.Prev() != nil || aKid.Next() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in AddKids: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("AddKids(K) can't cos K has siblings")
		}
	    if aKid.Parent() != nil && aKid.Parent() != p {
		fmt.Fprintf(os.Stdout, "FATAL in AddKids: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("E.AddKids(K) can't cos K has non-P parent")
		}
	    // All clear! Go ahead and add the kid.
	    _ = p.AddKid(aKid)
	}
	// FIXME JEEZ
	// For each new, set parent, then append whole mess to kids
	/*
	var FK = p.firstKid
	var LK = p.lastKid
	// Set the level now
	aKid.setLevel(p.Level() + 1)
	// Is the new kid an only kid ?
	if FK == nil && LK == nil {
		p.firstKid, p.lastKid = aKid, aKid
		aKid.SetParent(p)
		aKid.SetPrev(nil)
		aKid.SetNext(nil)
		return aKid
	}
	if !(FK != nil && LK != nil) {
		panic("BAD KID LINKS")
	}
	// So, replace the last kid
	if LK != nil {
		if LK.Parent() != p {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid dusnt know E")
		}
		if LK.Next() != nil {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid has a next kid")
		}
		LK.SetNext(aKid) // LK.nextKid = aKid
		aKid.SetPrev(LK) // aKid.prevKid = LK
		p.lastKid = aKid
		aKid.SetParent(p)
		return aKid
	}

	}
	fmt.Fprintf(os.Stdout, "FATAL in AddKids: " +
		"E<< %+v >> K<< %+v >>\n", p, rKids)
	panic("AddKids: Chaos!")
	*/
	return p
}

// FirstKid provides read-only access for other packages. Can return nil.
func (p *Linkset[T]) FirstKid() T {
	// return p.firstKid
	return p.kids[0]
}

// LastKid provides read-only access for other packages. Can return nil.
func (p *Linkset[T]) LastKid() T {
	// return p.lastKid
	ln := len(p.kids) 
	return p.kids[ln-1] 
}

// Prev provides read-only access for other packages. Can return nil.
// 2026.05: It needs access to the slice it is in, which means access
// to the parent, so if parent is nil, this returns nil. 
func (p *Linkset[T]) Prev() T {
	// return p.prevPeer
	if p.prnt == nil { return nil } 
	// Find in list and return preceding
	var ppk []*Nork = p.prnt.kids
	if ppk == nil || len(ppk) == 0 { return nil }
	idx := slices.Index(ppk, p)
	if idx == -1 || idx == 0 { return nil }
	return ppk[idx-1]
}

// Next provides read-only access for other packages. Can return nil.
// 2026.05: It needs access to the slice it is in, which means access
// to the parent, so if parent is nil, this returns nil. 
func (p *Linkset[T]) Next() T {
	// return p.nextPeer
	if p.prnt == nil { return nil } 
	// Find	in list	and return succeeding
	var ppk []*Nork = p.prnt.kids
	if ppk == nil || len(ppk) == 0 { return nil }
	idx := slices.Index(ppk, p)
	if idx == -1 || idx == len(p.prnt.kids)-1 { return nil }
	return ppk[idx+1]
}

func (p *Linkset[T]) KidsAsSlice() []T {
     	return p.kids
}

