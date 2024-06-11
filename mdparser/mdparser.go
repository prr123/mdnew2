// parser for markdown files
package mdparser


import (
	"fmt"
)

type MdNode struct {
	ch []*MdNode
	par *MdNode
	typ string
	blkSt int
	blkEnd int
	txtSt int
	txt []byte
//	att *Attribute
//	prop interface
}

type MdParser struct {
	blkMap map[byte]func(ps *MdPState) *MdNode
	lines []RLine
	max int
}

type MdPState struct {
	Doc *MdNode //top
	Node *MdNode // current parent node
	Blk *MdNode // prev block
	closed bool
	plin int
}

type RLine struct {
	linSt int
	linEnd int
	lintxt []byte
	eolChar int
}

func IsAlpha(let byte)(res bool) {
    res = false
    if (let >= 'a' && let <= 'z') || (let >= 'A' && let <= 'Z') { res = true}
    return res
}

func CleanRet (inp *[]byte) {

	ptr := 0
	for i:=0; i< len(*inp); i++ {
		if (*inp)[i] == '\r' {
			if (*inp)[i+1] != '\n' {
				(*inp)[ptr] = '\n'
				ptr++
			}
		} else {
			(*inp)[ptr] = (*inp)[i]
			ptr++
		}
	}
}


func InitParser(inp []byte) (p MdParser) {

	p.blkMap = make(map[byte]func(ps *MdPState) *MdNode)

	p.blkMap['#'] = p.ParseHeading
	p.blkMap['p'] = p.ParsePar
	p.blkMap[' '] = p.ParseEL
	p.blkMap['`'] = p.ParseCode
	p.blkMap['-'] = p.ParseUL
	p.blkMap['+'] = p.ParseUL
	p.blkMap['*'] = p.ParseUL
	p.blkMap['n'] = p.ParseOL
	p.blkMap['>'] = p.ParseQuote

	p.lines = GetLines(inp)
	p.max = len(p.lines)
	return p
}

func InitParseState(inp []byte) (pstate *MdPState) {

	var mdDoc MdNode
	mdDoc.blkSt = 0
	mdDoc.blkEnd = len(inp)
	mdDoc.typ = "doc"
	mdDoc.ch = nil
	mdDoc.par = nil

	var ps MdPState
	ps.Doc = &mdDoc
	ps.Blk = nil
	ps.Node = &mdDoc
	ps.closed = true

	return &ps
}

func (p *MdParser)ParseCode(ps *MdPState) *MdNode {
	fmt.Printf("parsing code\n")
	return nil
}

func  (p *MdParser)ParseQuote(ps *MdPState) *MdNode {
	fmt.Printf("parsing code\n")
	return nil
}

/*
type MdNode struct {
	ch []*MdNode
	par *MdNode
	typ string
	blkSt int
	blkEnd int
	txtSt int
	txt []byte
//	att *Attribute
//	prop interface
}
*/

func  (p *MdParser)ParseUL(ps *MdPState) *MdNode {
	fmt.Printf("parsing UL\n")

	l := p.lines[ps.plin]
	state:=0

	blk := &MdNode{}

	if ps.Node.typ != "UL" {
		blk.typ = "UL"
		blk.par = ps.Node
		blk.blkSt= l.linSt
		blk.blkEnd = -1
		ps.Node = blk
	} else {
		blk = ps.Node
	}

	liblk := &MdNode{
			typ: "LI",
			par: blk,
			blkSt: l.linSt,
			blkEnd: -1,
		}

//	suc := false
	loop := true
	for i:=1; i<len(l.lintxt); i++ {
		let := l.lintxt[i]
		switch state {
		case 0:
			if let == ' ' {state = 1}
		case 1:
			if let == ' ' {break}
			if let != ' ' {
				state = 2
				liblk.txtSt = i
				liblk.txt = l.lintxt[i:]
//				suc = true
				break
			}
		case 2:
			loop = false
		default:
			return nil
		}
		if !loop {break}
	}
	blk.ch = append(blk.ch,liblk)
	return blk
}

func (p *MdParser)ParseOL(ps *MdPState) *MdNode {
	fmt.Printf("parsing OL\n")

	return nil
}

func (p *MdParser)ParseHeading(ps *MdPState) *MdNode{
	l := p.lines[ps.plin]
	fmt.Printf("parsing heading: %s\n", string(l.lintxt))

	hdlev :=0
	state:=0
	txtst:=-1

	fin := false
	for i:=0; i<len(l.lintxt); i++ {

		let := l.lintxt[i]
		switch state {
		case 0:
			if let == '#' {hdlev++}
			if let == ' ' {state = 1}
		case 1:
			if let != ' ' {
				state = 2
				txtst = i
			}
		case 2:
			fin = true
		default:
		}
		if fin {break}
	}

	head := fmt.Sprintf("h%d",hdlev)
	txtSt := l.linSt + txtst
	blk := MdNode{
		typ: head,
		par: ps.Node,
		blkSt: l.linSt,
		blkEnd: l.linEnd,
		txtSt: txtSt,
		txt: l.lintxt[txtst:],
	}
	return &blk
}

func (p *MdParser)ParseEL(ps *MdPState) *MdNode{
	fmt.Println("parsing empty line")
	l := p.lines[ps.plin]
    blk := MdNode{
        typ: "br",
        par: ps.Node,
        blkSt: l.linSt,
        blkEnd: l.linEnd,
    }

	for i:=0; i<len(l.lintxt); i++ {
		let := l.lintxt[i]
		if let != ' ' {
fmt.Printf("not a empty line: %q\n",let)
			return nil
		}
	}
//fmt.Println("adding br")
//	p.Node.ch = append(p.Node.ch, &blk)
//	p.Blk = nil
	ps.closed = true
	return &blk
}

func (p *MdParser)ParsePar(ps *MdPState) *MdNode{
	fmt.Println("parsing paragraph")
//fmt.Printf("ps.close: %t\n%v\n",ps.closed, ps)

	l := p.lines[ps.plin]
	eoBlk:= false
	if l.lintxt[len(l.lintxt)-1]== ' ' && l.lintxt[len(l.lintxt)-2] == ' ' {
//fmt.Println("end of par 2ws")
		l.lintxt  =  l.lintxt[:len(l.lintxt)-2]
		eoBlk = true
	}

	blk := &MdNode{}
	if ps.closed {
		blk = &MdNode{
				typ: "p",
				par: ps.Node,
				blkSt: l.linSt,
				blkEnd: l.linEnd,
				txtSt: l.linSt,
				txt: l.lintxt,
			}
	} else {
		blk = ps.Blk
		blk.blkEnd = l.linEnd
		blk.txt = append(ps.Blk.txt, ' ')
		blk.txt = append(ps.Blk.txt, l.lintxt...)
	}

//fmt.Printf("par p return: %v\np.Blk:%v\n", p, p.Blk)
	ps.closed = false
	if eoBlk {ps.closed = true}
	return blk
}

func CloseBlk(ps *MdPState) {
	fmt.Println("closing block")

	if ps.Blk != nil {
		ps.Node.ch = append(ps.Node.ch, ps.Blk)
		ps.Blk = nil
	}
	return
}

func GetLines (inp []byte) (linList []RLine){

	linSt:=0
	linList = make([]RLine,0,128)

	for i:=0; i< len(inp); i++ {
		switch inp[i] {
		case '\n':
			newLine := RLine {
				linSt: linSt,
				linEnd: i,
				lintxt: inp[linSt:i],
				eolChar: 0,
			}
			if linSt == i  {newLine.eolChar = 1}
			if i-linSt >2 {
				if inp[i-2] == ' ' && inp[i-1] == ' ' {newLine.eolChar = 2}
			}
			linList = append(linList,newLine)
			linSt = i+1
		default:

		}
	}
	return linList
}


func (p *MdParser)Parse (ps *MdPState) (err error){

	linList := p.lines

	for i:=0; i< len(linList); i++ {
		line := linList[i]
//fmt.Printf("*** Line[%d]: \"%s\"\n", i+1, string(line.lintxt))
		ps.plin = i
		res:=&MdNode{}
		psold := ps.closed
		tmp := ""
		if len(line.lintxt) == 0 {
			res = p.ParseEL(ps)
			tmp = "el"
		} else {
			flet := line.lintxt[0]
			plet := flet
			if IsAlpha(flet) {plet='p'}
			f, ok := p.blkMap[plet]
			if !ok {
				fmt.Printf("error -- line[%d] first letter unknown: %q\n", i, plet)
				continue
			}
//fmt.Printf("before parse:%q %t\n%v\n", plet, ps.closed,ps)
			res=f(ps)
			tmp = string(plet)
		}

		tmp = "res: " + tmp
		PrintNode(res, tmp)
//fmt.Printf("after parse:%q %t\n%v\n", plet, ps.closed,ps)

/*
		if ps.Blk != nil {
			if ps.Blk.typ != res.typ && !psold{
				ps.Node.ch = append(ps.Node.ch, ps.Blk)
				ps.closed = true
			}
		}
*/
		if ps.closed && !psold {
				ps.Node.ch = append(ps.Node.ch, ps.Blk)
		}
		if ps.closed {
			ps.Node.ch = append(ps.Node.ch, res)
		}
		ps.Blk = res

//fmt.Printf("after res parse %q: %v\np.Blk:%v\n",plet, p, p.Blk)
//fmt.Printf("[%d]: %t %q\"%s\"\n", i+1, res, plet, string(line.lintxt))
	}

	return nil
}


func PrintNode(n *MdNode, title string) {

	fmt.Printf("\n******** Node %s ***********\n", title)
	if n == nil {
		fmt.Println("no node")
		fmt.Printf("****** End Node %s *********\n\n", title)
		return
	}
	fmt.Printf("Typ: %s\n", n.typ)
 	fmt.Printf("st: %d end: %d\n", n.blkSt, n.blkEnd)
	fmt.Printf("children: %d\n", len(n.ch))
	if n.par == nil {
		fmt.Printf("parent: none\n")
	} else {
		fmt.Printf("parent: %s\n", n.par.typ)
	}
	fmt.Printf("txt: %s\n", n.txt)

//	if par == nil {return}
	fmt.Printf("Children [%d]\n", len(n.ch))
	if len(n.ch) == 0 {
		fmt.Printf("****** End Node %s *********\n\n", title)
		return
	}
	for i:= 0; i< len(n.ch); i++ {
		cNode := n.ch[i]
		str := fmt.Sprintf("child: %d", i +1)
//fmt.Printf("** %s **\n", str)
		PrintNode(cNode, str)
	}

	fmt.Printf("****** End Node %s *********\n\n", title)

}
