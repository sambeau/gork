package main

import (
	"fmt"
	"io"
)

const indent = "    "

// - Texts - - - - - - - - - -

// All node types implement the Lexed interface.
type Lexed interface {
	Line() int
	Column() int
}

type Pos struct {
	LineNo   int
	ColumnNo int
}

func (p Pos) Line() int {
	return p.LineNo
}

func (p Pos) Column() int {
	return p.ColumnNo
}

// type Token struct {
// 	Lexed
// 	Text string
// }

// - Texts - - - - - - - - - -

type Text struct {
	Lexed
	StringV string
}

func (txt Text) String() string {
	return txt.StringV
}

// - Texts - - - - - - - - - -

type Handle struct {
	Lexed
	StringV string
}

func (h Handle) String() string {
	return "{" + h.StringV + "}"
}

// - - - - - - - - - - -

type Texter interface {
	String() string
}

type Texts struct {
	List []Texter
}

func (d Texts) Empty() bool {
	return len(d.List) == 0
}

func (d *Texts) Add(text Texter) {
	d.List = append(d.List, text)
}

func (d *Texts) Merge(d1 Texts) {
	d.List = append(d.List, d1.List...)
}

func (d Texts) String() string {
	s := ""
	nd := false
	for _, texts := range d.List {
		if nd {
			s += " "
		}
		s += texts.String()
		nd = true
	}
	return s
}
func (d Texts) Evaluate() Texter {
	return Texter(d)
}

// - Expr - - - - - - - - - -

type Evaluatable interface {
	Texter
	Evaluate() Texter
}

type Block struct {
	Pos
	Nodes []Evaluatable
}

func (x *Block) Add(n Evaluatable) {
	x.Nodes = append(x.Nodes, n)
}

func (x Block) Evaluate() Texter {
	val := Texter(Texts{})
	for _, e := range x.Nodes {
		val = e.Evaluate()
	}
	return val
}

func (x Block) String() string {
	s := "("
	nd := false
	for _, node := range x.Nodes {
		if nd {
			s += " "
		}
		s += node.String()
		nd = true
	}
	s += ")"
	return s
}

type Node struct {
	Pos
	Node Evaluatable
}

func (n Node) String() string {
	return n.Node.String()
}

func (n Node) Evaluate() Texter {
	return n.Node.Evaluate()
}

type IfNode struct {
	Pos
	Cond Node
	Then Block
	Else Block
}

func (n IfNode) String() string {
	return "if " + n.Cond.String() + " then " + n.Then.String() + " else " + n.Else.String()
}

func (n IfNode) Evaluate() Texter {
	return n.Cond.Evaluate()
}

type IsCond struct {
	Pos
	Name  string
	Truth bool
}

func (n IsCond) String() string {
	if !n.Truth {
		return "is not " + n.Name
	} else {
		return "is " + n.Name
	}
}

func (n IsCond) Evaluate() Texter {
	return Text{n.Pos, n.Name}
}

// - Traits - - - - - - - - - -

type Traits map[string]bool

func (ts1 Traits) Merge(ts2 Traits) {
	for k, v := range ts2 {
		ts1[k] = v
	}
}

func (t Traits) Pprint(w io.Writer, i string) {
	fmt.Fprintf(w, "%sis", i)
	nd := false
	for k, v := range t {
		if nd {
			fmt.Fprintf(w, ",")
		}
		if v == false {
			fmt.Fprintf(w, " not")
		}
		fmt.Fprintf(w, " %s", k)
		nd = true
	}
	fmt.Fprintf(w, "%s\n", i)
}

// - Location - - - - - - - - - -

type Location struct {
	Name        string
	Description Texts
	Traits      Traits
}

func (l Location) Pprint(w io.Writer, i string) {
	i2 := i + indent
	fmt.Fprintf(w, "%slocation %s [\n", i, l.Name)
	fmt.Fprintf(w, "%s%s\n", i2, l.Description.String())
	l.Traits.Pprint(w, i2)
	fmt.Fprintf(w, "%s]\n", i)
}

// - Game - - - - - - - - - -

type Game struct {
	Name        string
	Title       Texts
	By          Texts
	Description Texts
	Version     Texts
	Date        Texts
	Traits      Traits
	Locations   map[string]Location
	Current     string
}

func (g Game) Pprint(w io.Writer, i string) {
	i2 := i + indent
	fmt.Fprintf(w, "%sgame [\n", i)
	fmt.Fprintf(w, "%stitle %s\n", i2, g.Title.String())
	if !g.By.Empty() {
		fmt.Fprintf(w, "%sby %s\n", i2, g.By.String())
	}
	fmt.Fprintf(w, "%s\n", i2)
	if !g.Version.Empty() {
		fmt.Fprintf(w, "%sversion %s\n", i2, g.Version.String())
	}
	if !g.Date.Empty() {
		fmt.Fprintf(w, "%sdate %s\n", i2, g.Date.String())
	}
	if !g.Version.Empty() || !g.Date.Empty() {
		fmt.Fprintf(w, "%s\n", i2)
	}
	fmt.Fprintf(w, "%s%s\n", i2, g.Description.String())
	fmt.Fprintf(w, "%s\n", i2)
	fmt.Fprintf(w, "%sstart %s\n", i2, g.Current)
	fmt.Fprintf(w, "%s\n", i2)
	nd := false
	for _, v := range g.Locations {
		if nd {
			fmt.Fprintf(w, "\n")
		}
		v.Pprint(w, i2)
		nd = true
	}
	fmt.Fprintf(w, "%s]\n", i)
}
