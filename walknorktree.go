package nork

import(
//  "io"
    TL "github.com/fbaube/treelinks"
)

type InspectorFunc func(pNode *Nork) error
type  StringerFunc func(pNode *Nork) string 

/*

func (p *Nork) StringserTree(f StringerFunc, w io.Writer) error {
    	s := f(p) + "\n" 
	// Write(p []byte) (n int, err error) 
	_, _ = w.Write([]byte(s+"\n")) // e hidn
	pKid := p.Links.FirstKid()
	for pKid != nil {
	    	s = f((pKid.Links)(*Nork))+"\n"
		_, _ = w.Write([]byte(s+"\n")) // e hidn
/* FIXME		if e = pKid.StringserTree(f, w); e != nil {
			return e
		} * /
		pKid = pKid.Next()
	}
	return nil
}

*/

func (p *Nork) FirstKid() TL.Treelinkser {
      return p.Links.FirstKid()
}
      
/*

func DumpTree(p *Nork, f StringerFunc, w io.Writer) error {
	var e error
    	s := p.f()
	w.Writeln(s)
	pKid := p.FirstKid()
	for pKid != nil {
	    	s := p.f()
		w.Writeln(s)
		if e = DumpTree(pKid, f, w); e != nil {
			return e
		}
		pKid = pKid.NextPeer()
	}
	return nil
}

/*

// func InspectTree used to be func WalkNorders 
// .
func InspectTree(p *Nork, f InspectorFunc) error {
	var e error
	if e = f(p); e != nil {
		return e
	}
	pKid := p.FirstKid()
	for pKid != nil {
		if e = InspectTree(pKid, f); e != nil {
			return e
		}
		pKid = pKid.NextPeer()
	}
	return nil
}

func InspectTreeWithPreAndPost(p *Nork,
	f0 InspectorFunc, f1 InspectorFunc) error {

	var e error
	// PRE
	if e = f0(p); e != nil {
		return e
	}
	// KIDS
	pKid := p.FirstKid()
	for pKid != nil {
		if e = InspectTreeWithPreAndPost(pKid, f0, f1); e != nil {
			return e
		}
		pKid = pKid.NextPeer()
	}
	// POST
	if e = f1(p); e != nil {
		return e
	}
	return nil
}

*/

