package nork

/*
import (
	"fmt"
	"os"
	"slices"
)
*/

// func (n *Node[D]) AddChild(child *Node[D]) {
//     n.Children = append(n.Children, child)
// }

// THESE FUNC'S ARE HACKY !!
// They use "Kids" like they are still a link list,
// but now in this rewrite (2026.04) they are a slice.

// HasKids is duh.
// func (p *Nork) [D any] HasKids() bool {
func (p *Linkset[T]) HasKids() bool {
	// return p.FirstKid() != nil && p.LastKid() != nil
	return p.kids != nil && len(p.kids) != 0
}

// Parent returns the parent, duh.
// func (p *Nork)(p *Linkset[T]) Parent() T {
func (p *Linkset[T]) Parent() T {
	return p.prnt
}

// SetParent has no side effects.
func (p *Linkset[T]) SetParent(p2 T) {
	p.prnt = p2
}

// SetPrevPeer has no side effects.
func (p *Linkset[T]) SetPrevPeer(p2 T) {
	p.SetPrevPeer(p2)
}

// SetNextPeer has no side effects.
func (p *Linkset[T]) SetNextPeer(p2 T) {
	p.SetNextPeer(p2)
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

* / FIXME

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
// NOTE that just about ALL old (linked list) code is ripped out. 
func (p *Linkset[T]) AddKid(aKid *Nork) T { // returns aKid
	// fmt.Printf("nord: ptrs? aKid<%T> p<%T> \n", aKid, p)
	if aKid.PrevPeer() != nil || aKid.NextPeer() != nil {
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
	    if aKid.PrevPeer() != nil || aKid.NextPeer() != nil {
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
		aKid.SetPrevPeer(nil)
		aKid.SetNextPeer(nil)
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
		if LK.NextPeer() != nil {
			fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
			panic("E.AddKid: E's last kid has a next kid")
		}
		LK.SetNextPeer(aKid) // LK.nextKid = aKid
		aKid.SetPrevPeer(LK) // aKid.prevKid = LK
		p.lastKid = aKid
		aKid.SetParent(p)
		return aKid
	}

	}
	fmt.Fprintf(os.Stdout, "FATAL in AddKids: " +
		"E<< %+v >> K<< %+v >>\n", p, rKids)
	panic("AddKids: Chaos!")
	* / FIXME
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

// PrevPeer provides read-only access for other packages. Can return nil.
// 2026.05: It needs access to the slice it is in, which means access
// to the parent, so if parent is nil, this returns nil. 
func (p *Linkset[T]) PrevPeer() T {
	// return p.prevPeer
	if p.prnt == nil { return nil } 
	// Find in list and return preceding
	var ppk []*Nork = p.prnt.kids
	if ppk == nil || len(ppk) == 0 { return nil }
	idx := slices.Index(ppk, p)
	if idx == -1 || idx == 0 { return nil }
	return ppk[idx-1]
}

// NextPeer provides read-only access for other packages. Can return nil.
// 2026.05: It needs access to the slice it is in, which means access
// to the parent, so if parent is nil, this returns nil. 
func (p *Linkset[T]) NextPeer() T {
	// return p.nextPeer
	if p.prnt == nil { return nil } 
	// Find	in list	and return succeeding
	var ppk []*Nork = p.prnt.kids
	if ppk == nil || len(ppk) == 0 { return nil }
	idx := slices.Index(ppk, p)
	if idx == -1 || idx == len(p.prnt.kids)-1 { return nil }
	return ppk[idx+1]
}

func (p *Linkset[T]) KidsAsSlice() []*Nork {
     	return p.kids
}

*/
