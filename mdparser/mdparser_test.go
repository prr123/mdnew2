package mdparser

import (
	"testing"
	"os"
	"fmt"
)

func TestPrLines(t *testing.T) {

	tstFilnam := "/home/peter/go/src/goDemo/mdnew2/mdFiles/testULNest1.md"
	content, err := os.ReadFile(tstFilnam)
	if err != nil {
		t.Error("cannot read test file!")
		return
	}
	p := InitParser(content)
	lines := p.lines

	PrLines(lines)
}

func TestLines(t *testing.T) {
	tstFilnam := "/home/peter/go/src/goDemo/mdnew2/mdFiles/test3A1.md"
	content, err := os.ReadFile(tstFilnam)
	if err != nil {
		t.Error("cannot read test file!")
		return
	}

	p := InitParser(content)
	ps := InitParseState(content)
//	lines, err := GetLines(content, MdP)
//	if err != nil {t.Error("cannot get Lines!")}

	fmt.Println("****** raw text lines *******")
	lines := p.lines
	for i:=0; i<len(lines);i++ {
		linst := lines[i].linSt
		linend:= lines[i].linEnd
		fmt.Printf("[%d]: %s\n", i+1, string(content[linst:linend]))
	}
	fmt.Println("**** end raw text lines *****")
	fmt.Println()

	err = p.Parse(ps)
	if err != nil {t.Error("cannot parse Lines!")}
	PrintNode(ps.Doc, "doc lines")
}

func TestPar(t *testing.T) {
	tstFilnam := "/home/peter/go/src/goDemo/mdnew2/mdFiles/testPar.md"
	content, err := os.ReadFile(tstFilnam)
	if err != nil {
		t.Error("cannot read testPar file!")
		return
	}

	p := InitParser(content)
	ps := InitParseState(content)

	fmt.Println("****** raw text lines *******")
	lines := p.lines
	for i:=0; i<len(lines);i++ {
		linst := lines[i].linSt
		linend:= lines[i].linEnd
		fmt.Printf("[%d]: %s\n", i+1, string(content[linst:linend]))
	}
	fmt.Println("**** end raw text lines *****")
	fmt.Println()

	err = p.Parse(ps)
	if err != nil {t.Error("cannot parse Lines!")}
	PrintNode(ps.Doc, "doc paragraphs")
}


func TestHeadings(t *testing.T) {
	tstFilnam := "/home/peter/go/src/goDemo/mdnew2/mdFiles/testHeadings.md"
	content, err := os.ReadFile(tstFilnam)
	if err != nil {
		t.Error("cannot read test file!")
		return
	}

	p := InitParser(content)
	ps := InitParseState(content)

	fmt.Println("****** raw text lines *******")
	lines := p.lines
	for i:=0; i<len(lines);i++ {
		linst := lines[i].linSt
		linend:= lines[i].linEnd
		fmt.Printf("[%d]: %s\n", i+1, string(content[linst:linend]))
	}
	fmt.Println("**** end raw text lines *****")
	fmt.Println()

	err = p.Parse(ps)
	if err != nil {t.Error("cannot parse Lines!")}
	PrintNode(ps.Doc, "doc headings")
}

func TestUL(t *testing.T) {
	tstFilnam := "/home/peter/go/src/goDemo/mdnew2/mdFiles/testUL.md"
	content, err := os.ReadFile(tstFilnam)
	if err != nil {t.Error("cannot read test file!")}

	p := InitParser(content)
	ps := InitParseState(content)

	fmt.Println("****** raw text lines *******")
	lines := p.lines
	for i:=0; i<len(lines);i++ {
		linst := lines[i].linSt
		linend:= lines[i].linEnd
		fmt.Printf("[%d]: %s\n", i+1, string(content[linst:linend]))
	}
	fmt.Println("**** end raw text lines *****")
	fmt.Println()

	err = p.Parse(ps)
	if err != nil {t.Error("cannot parse nd!")}
	PrintNode(ps.Doc, "doc ul")
}

