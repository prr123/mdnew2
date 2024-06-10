package main

// example for https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"time"

    "github.com/goccy/go-yaml"
//	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
//	"github.com/gomarkdown/markdown/html"
	mdjs "goDemo/md/mdjs"
	"github.com/gomarkdown/markdown/parser"
	util "github.com/prr123/utility/utilLib"
)

func RenderDom(doc ast.Node, renderer mdjs.Renderer) []byte {
    var buf bytes.Buffer
    renderer.RenderHeader(&buf, doc)
    ast.WalkFunc(doc, func(node ast.Node, entering bool)(walk ast.WalkStatus) {
        xy := renderer.RenderNode(&buf, node, entering)
//        fmt.Printf("walk status: %d %s\n", xy, buf)
        return xy
    })
    renderer.RenderFooter(&buf, doc)
    return buf.Bytes()
}


var mds = `# header

Sample text.

[link](http://example.com)
`

//var printAst = true


func mdToJsDom(md []byte, dbg bool) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	if dbg {
		fmt.Print("--- AST tree:\n")
		ast.Print(os.Stdout, doc)
		fmt.Print("\n")
	}

	// create HTML renderer with extensions
//	htmlFlags := html.CommonFlags | html.HrefTargetBlank
//	opts := html.RendererOptions{Flags: htmlFlags}
//	renderer := html.NewRenderer(opts)

	renderer := mdjs.NewRenderer()

	return RenderDom(doc,*renderer)
}

func main() {

    numArgs := len(os.Args)
	var md []byte
	var mds =`
# header

Sample text.

[link](http://example.com)
`

	md = []byte(mds)

	flags:=[]string{"dbg","md", "script"}

    useStr := "[/md=<markdown file>] [/script=<script file>] [/dbg]"
    helpStr := fmt.Sprintf("help: The program cmvert md files to a js Dom script\n")

    if numArgs > len(flags)+1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }


    if numArgs == 2 {
        if os.Args[1] == "help" {
            fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
            fmt.Printf("%s\n", helpStr)
            os.Exit(1)
        }
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg := false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    mdFilnam := "test"
	mdFullFilnam := ""
    mdval, ok := flagMap["md"]
    if ok {
        if mdval.(string) == "none" {log.Fatalf("error -- no markdown file provided with /md flag!")}
        mdFilnam = mdval.(string)
//	idx := bytes.Index(mdFilnam,".md")
		mdFullFilnam = "mdFiles/" + mdFilnam + ".md"
		md, err = os.ReadFile(mdFullFilnam)
		if err != nil {log.Fatalf("error -- cannot read md: %v", err)}
	}

    outFilnam := ""
    oval, ok := flagMap["script"]
    if ok {
        if oval.(string) == "none" {log.Fatalf("error -- no script file provided with /script flag!")}
        outFilnam = oval.(string)
//	idx := bytes.Index(mdFilnam,".md")
		outFilnam = "script/" + outFilnam + ".js"
	} else {
		outFilnam = "script/" + mdFilnam + ".js"
	}

	outfil, err := os.Create(outFilnam)
	if err != nil {log.Fatalf("error -- cannot create script file: %v", err)}
	defer outfil.Close()

	// check whether md contains a yaml section
	var yamlSec []byte

	mdSt := 0
	if idx := bytes.Index(md, []byte("====\n")); idx > -1 {
		idx2 := bytes.Index(md[5:], []byte("====\n"))
		if idx2 > -1 {
			yamlSec = md[5:idx2+5]
			mdSt = idx2+11
		}

	}


    if dbg {
		if len(mdFilnam)  == 0 {
			log.Printf("debug -- no md file!\n")
		} else {
        	log.Printf("debug -- md file: %s\n",mdFilnam)
		}
		log.Printf("debug -- script file: %s\n",outFilnam)
		log.Printf("debug -- yaml section: %d\n", len(yamlSec))
    }

	var DHead mdjs.DocHead
	if len(yamlSec)>0 {
		fmt.Printf("dbg -- yaml:\n\"%s\"\n", string(yamlSec))
		err = yaml.Unmarshal(yamlSec, &DHead)
    	if err != nil {log.Fatalf("error -- yaml unmarshal: %v\n", err)}
		fmt.Printf("DHead: %v\n", DHead)
		layout := "2/1/2006"
		DHead.Date, err = time.Parse(layout, DHead.DateStr)
    	if err != nil {log.Fatalf("error -- time conversion: %v\n", err)}
		if dbg {PrintDH(DHead)}
	}

	script := mdToJsDom(md[mdSt:], dbg)

//	fmt.Printf("\n\n--- Markdown:\n%s\n\n--- jsDom:\n%s\n", md, script)

	if outfil != nil {
		_, err = outfil.Write(script)
		if err !=nil {log.Fatalf("error -- writing to script file! %v\n",err)}
	}
}

/*
    Title string `yaml:"title"`
    Author string `yaml:"author"`
    Date time.Time `yaml:"date"`
    Summary string `yaml:"summary"`
*/
func PrintDH(dhead mdjs.DocHead) {

	fmt.Println("****** Doc Head *******")
	fmt.Printf("Title:  %s\n", dhead.Title)
	fmt.Printf("Author: %s\n", dhead.Author)
	fmt.Printf("Summary:\n\"%s\"\n", dhead.Summary)
	fmt.Printf("DateStr: %s\n", dhead.DateStr)
	fmt.Printf("Date: %s\n", dhead.Date.Format("2/1/2006"))


	fmt.Println("**** End Doc Head *****")

}
