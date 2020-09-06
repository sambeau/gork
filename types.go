package main

import (
	"fmt"
	"io"
)

const indent = "    "

// - Descriptions - - - - - - - - - -

// All node types implement the Node interface.
type Node interface {
	Line() int
	Column() int
}

type Pos struct {
	line   int
	column int
}

func (p Pos) Line() int {
	return p.line
}

func (p Pos) Column() int {
	return p.column
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
}

type Token struct {
	Node
	Text string
}

// - Descriptions - - - - - - - - - -

type Text struct {
	Node
	String string
}

func (txt Text) Describe() string {
	return string(txt.String)
}

// - - - - - - - - - - -

type Describable interface {
	Describe() string
}

type Descriptions struct {
	list []Describable
}

func (d Descriptions) Empty() bool {
	return len(d.list) == 0
}

func (d *Descriptions) Add(describable Describable) {
	d.list = append(d.list, describable)
}

func (d *Descriptions) Merge(d1 Descriptions) {
	d.list = append(d.list, d1.list...)
}

func (d Descriptions) Describe() string {
	s := ""
	nd := false
	for _, describable := range d.list {
		if nd {
			s += " "
		}
		s += describable.Describe()
		nd = true
	}
	return s
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
	Description Descriptions
	Traits      Traits
}

func (l Location) Pprint(w io.Writer, i string) {
	i2 := i + indent
	fmt.Fprintf(w, "%slocation %s [\n", i, l.Name)
	fmt.Fprintf(w, "%s\"%s\"\n", i2, l.Description.Describe())
	l.Traits.Pprint(w, i2)
	fmt.Fprintf(w, "%s]\n", i)
}

// - Game - - - - - - - - - -

type Game struct {
	Title       Descriptions
	By          Descriptions
	Description Descriptions
	Version     Descriptions
	Date        Descriptions
	Locations   map[string]Location
	Current     string
}

func (g Game) Pprint(w io.Writer, i string) {
	i2 := i + indent
	fmt.Fprintf(w, "%sgame [\n", i)
	fmt.Fprintf(w, "%stitle \"%s\"\n", i2, g.Title.Describe())
	if !g.By.Empty() {
		fmt.Fprintf(w, "%sby \"%s\"\n", i2, g.By.Describe())
	}
	fmt.Fprintf(w, "%s\n", i2)
	if !g.Version.Empty() {
		fmt.Fprintf(w, "%sversion %s\n", i2, g.Version.Describe())
	}
	if !g.Date.Empty() {
		fmt.Fprintf(w, "%sdate \"%s\"\n", i2, g.Date.Describe())
	}
	if !g.Version.Empty() || !g.Date.Empty() {
		fmt.Fprintf(w, "%s\n", i2)
	}
	fmt.Fprintf(w, "%s\"%s\"\n", i2, g.Description.Describe())
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
