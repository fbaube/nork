package nork

import (
)

func (p *Cnty) AbsFP() string { return p.absPath }
func (p *Cnty) RelFP() string { return p.relPath }
func (p *Cnty) SetAbsFP(s string) { p.absPath = s }
func (p *Cnty) SetRelFP(s string) { p.relPath = s }

/* REWRITE to add GNORK func's 

// func (n *Node[D]) AddChild(child *Node[D]) {
//     n.Children = append(n.Children, child)
// }

// HasKids is duh.
// func (p *Nork) [D any] HasKids() bool {
func (p *Nork) HasKids() bool {
	return p.FirstKid() != nil && p.LastKid() != nil
}

// Parent returns the parent, duh.
// func (p *Nork) Parent() Norker {
func (p *Nork) Parent() Norker {
	return p.Parent()
}

// SetParent has no side effects.
func (p *Nork) SetParent(p2 Norker) {
	p.prnt = p2.(*Nork)
}

// SetPrevPeer has no side effects.
func (p *Nork) SetPrevPeer(p2 Norker) {
	p.PrevPeer = p2
}

// SetNextPeer has no side effects.
func (p *Nork) SetNextPeer(p2 Norker) {
	p.nextPeer = p2
}

// SetFirstKid has no side effects.
func (p *Nork) SetFirstKid(p2 Norker) {
	p.firstKid = p2
}

// SetLastKid has no side effects.
func (p *Nork) SetLastKid(p2 Norker) {
	p.lastKid = p2
}

// AddKid adds the supplied node as the last kid, and returns
// it (i.e. the new last kid), now linked into the tree.
func (p *Nork) AddKid(aKid Norker) Norker { // returns aKid
	// fmt.Printf("nord: ptrs? aKid<%T> p<%T> \n", aKid, p)
	if aKid.PrevPeer() != nil || aKid.NextPeer() != nil {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("AddKid(K) can't cos K has siblings")
	}
	if aKid.Parent() != nil && aKid.Parent() != p {
		fmt.Fprintf(os.Stdout, "FATAL in AddKid: Tag<< %+v >> kid<< %+v >>\n", p, aKid)
		panic("E.AddKid(K) can't cos K has non-P parent")
	}
	var FK = p.FirstKid()
	var LK = p.LastKid()
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
	fmt.Fprintf(os.Stdout, "FATAL in AddKid: E<< %+v >> K<< %+v >>\n", p, aKid)
	panic("AddKid: Chaos!")
}

// AddKids adds the supplied nodes as kids, after any pre-existing
// kids, and returns the parent. 
func (p *Nork) AddKids(rKids []Norker) Norker { // returns p 
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
	* /
	return p
}

// FirstKid provides read-only access for other packages. Can return nil.
func (p *Nork) FirstKid() Norker {
	return p.firstKid
}

// LastKid provides read-only access for other packages. Can return nil.
func (p *Nork) LastKid() Norker {
	return p.lastKid
}

// PrevPeer provides read-only access for other packages. Can return nil.
func (p *Nork) PrevPeer() Norker {
	return p.prevKid
}

// NextPeer provides read-only access for other packages. Can return nil.
func (p *Nork) NextPeer() Norker {
	return p.nextKid
}

func (p *Nork) KidsAsSlice() []Norker {
	var pp []Norker
	c := p.FirstKid() // p.firstKid
	for c != nil {
		pp = append(pp, c)
		c = c.NextPeer() // c.nextKid
	}
	return pp
}

*/

