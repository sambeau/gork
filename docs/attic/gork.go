package main

import (
	"fmt"
)

type command struct {
	name     string
	synonyms []string
}

type object struct {
	id          string
	names       []string
	title       string
	description string
	commands    []command
	states      map[string]string
}

type prop struct {
	object
	portable bool
}

type actor struct {
	object
}

type location struct {
	object
	props  []*prop
	actors []*actor
}

func (l location) describe() {
	fmt.Printf("In %s\n", l.title)
	fmt.Printf("You are in %s.\n", l.description)
	for _, p := range l.props {
		fmt.Printf("There is %s here.\n", p.title)
	}
}

type inventory struct {
	items []prop
}

type game struct {
	currentLocation *location
	locations       map[string]*location
}

func (g game) move(loc string) {
	g.currentLocation = g.locations[loc]
	g.currentLocation.describe()
}

func (g game) start(loc string) {
	g.move(loc)
}

func main() {
	world := game{
		locations: map[string]*location{
			"hut": &location{
				object: object{
					id:          "hut",
					title:       "a hut",
					description: "a dingy hut",
				},
				props: []*prop{
					&prop{
						object: object{
							id:          "lamp",
							title:       "a lamp",
							description: "a small brass lamp",
						},
						portable: true,
					},
				},
			},
		},
	}

	world.start("hut")
}
