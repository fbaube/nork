package nork

import (
	"fmt"
	"io"
	// "os"
	S "strings"
)

var printTreeTo, printCssTreeTo io.Writer

// LinePrefixString provides indentation and
// should start a line of display/debug.
//
// It does not end the string with (white)space.
// .
func (p Nork) LinePrefixString() string {
	if p.IsRoot() { // && p.Parent == nil
		return "[R]"
	// } else if p.Level() == 0 && p.Parent() != nil {
	} else if p.Level() == 0 || p.Parent() != nil {
		// return fmt.Sprintf("[%d]", p.seqID)
		return "[?!R?!]"
	} else {
		// (spaces)[lvl:seq]"
		// func S.Repeat(s string, count int) string
		return fmt.Sprintf("%s[%02d]", // "%s[%02d:%02d]", 
			S.Repeat("  ", p.level-1), p.level) // ,p.seqID)
	}
}

func yn(b bool) string {
	if b {
		return "Y"
	} else {
		return "n"
	}
}

func (p *Nork) LineSummaryString() string {
	var sb S.Builder
	if p.IsRoot() {
		sb.WriteString("ROOT ")
	}
	/* more debugging
	if p.PrevPeer() != nil {	sb.WriteString("P ") }
	if p.Parent()   == nil {	sb.WriteString("NOPARENT ") }
	if p.NextPeer() != nil { sb.WriteString("N ") }
	if p.HasKids()  { sb.WriteString("kid(s) ") }
	*/
	if p.relPath == "" {
		sb.WriteString("NOPATH")
	} else {
		sb.WriteString(p.relPath)
	}
	return (sb.String())
}

func (p *Nork) PrintTree(w io.Writer) error {
	// println("PrintTree: could use printer fn")
	if w == nil {
		return nil
	}
	// printTreeTo = w
	// e := p.InfosAsTree(w)
	e := p.StringserTree(InfosG, w)
	if e != nil {
		println("nordPrintOneLiner ERR:", e.Error())
		return e
	}
	return nil
}

/*
func (p *Nork) PrintCssTree(w io.Writer) error {
	if w == nil {
		return nil
	}
	printTreeTo = w
	e := InspectTreeWithPreAndPost(p,
		nordPrintCssOneLinerPre, nordPrintCssOneLinerPost)
	if e != nil {
		println("nordPrintCssLine ret'd ERR:", e.Error())
		return e
	}
	return nil
}
*/

func (p *Nork) Echo() string {
	return "FIXME ECHO"
}
func (p *Nork) Infos() string {
	return "FIXME INFOS"
}
func (p *Nork) Debug() string {
	return "FIXME DEBUG"
}

func EchoG(p *Nork) string {
	return p.Echo()
}
func InfosG(p *Nork) string {
	return p.Infos()
}
func DebugG(p *Nork) string {
	return p.Debug()
}

func nordPrintCssOneLinerPre(p Norker) error {
	// firstEntry := true
	smry := p.Infos()

	if p.IsDir() {
		// if firstEntry {
		//  } else {
		// <li><details><summary><i>Ice</i> giants</summary>
		// <ul>
		fmt.Fprintf(printTreeTo, "<li><details><summary>"+
			smry+"</summary>\n<!--ul> ")
		// }
	} else {
		// Do both Pre AND Post
		fmt.Fprintf(printTreeTo, "<li>"+smry+"</li>\n")
	}
	fmt.Fprintf(printTreeTo, p.Infos())
	// firstEntry = false
	return nil
}

func nordPrintCssOneLinerPost(p Norker) error {
	if p.IsDir() {
		fmt.Fprintf(printTreeTo, "</li>\n")
	} // else {
	// Do both Pre AND Post
	// fmt.Fprintf(printTreeTo, "<li>"+smry+"</li>\n")
	// }
	fmt.Fprintf(printTreeTo, p.Infos())
	return nil
}
