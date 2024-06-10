package mdparser

import (
	"testing"
	"os"
	"fmt"
)


func TestLines(t* testing.T) {
	tstFilnam := "/home/peter/go/src/goDemo/mdnew/mdFiles/test3A1.md"
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
/*
func TestHeadings(t* testing.T) {
	tstFilnam := "/home/peter/go/src/goDemo/mdnew/mdFiles/testHeadings.md"
	content, err := os.ReadFile(tstFilnam)
	if err != nil {
		t.Error("cannot read test file!")
		return
	}

	MdP := InitParser(content)
	lines, err := GetLines(content)
	if err != nil {t.Error("cannot get Lines!")}

	for i:=0; i<len(lines);i++ {
		linst := lines[i].linSt
		linend:= lines[i].linEnd
		fmt.Printf("[%d]: %s\n", i+1, string(content[linst:linend]))
	}

	err = Parse(lines, MdP)
	if err != nil {t.Error("cannot parse Lines!")}
	PrintNode(MdP.Doc, "doc headings")
}

func TestUL(t* testing.T) {
	tstFilnam := "/home/peter/go/src/goDemo/mdnew/mdFiles/testUL.md"
	content, err := os.ReadFile(tstFilnam)
	if err != nil {t.Error("cannot read test file!")}

	MdP := InitParser(content)
	lines, err := GetLines(content)
	if err != nil {t.Error("cannot get Lines!")}

	for i:=0; i<len(lines);i++ {
		linst := lines[i].linSt
		linend:= lines[i].linEnd
		fmt.Printf("[%d]: %s\n", i+1, string(content[linst:linend]))
	}

	err = Parse(lines, MdP)
	if err != nil {t.Error("cannot parse nd!")}
	PrintNode(MdP.Doc, "doc ul")
}

*/
