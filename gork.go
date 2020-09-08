//gork.go

package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const version = "0.0.1"
const fileVersion = "0,0"
const identifier = "GORK!"
const ENDH = '\x0a'
const ENDC = '\x03'

func init() {
	gob.Register(Text{})
	gob.Register(Texts{})
	gob.Register(Pos{})
	gob.Register(Handle{})
	gob.Register(Block{})
	gob.Register(Node{})
	gob.Register(IfNode{})
	gob.Register(IsCond{})
}

func interleave(in string) string {
	n := 0
	c := ""
	out := ""
	for _, s := range in {
		switch n {
		case 0:
			c = "?"
		case 1:
			c = "%%"
		case 2:
			c = "??"
		case 3:
			c = "~"
		case 4:
			c = "=?"
		case 5:
			c = "???"
		case 6:
			c = "*"
		case 7:
			c = "?\\"
		}
		out += c + string(s)
		n += 1
		if n == 8 {
			n = 0
		}
	}
	return out
}

func saveGame(filename string, game Game) {
	var gamefile bytes.Buffer
	w := zlib.NewWriter(&gamefile)

	enc := gob.NewEncoder(w)
	err := enc.Encode(game)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	w.Flush()

	header := fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s%c",
		identifier,
		fileVersion,
		game.Name,
		game.Version.String(),
		game.By.String(),
		time.Now().Format(time.RFC3339),
		ENDH,
	)

	checksum := sha256.New()
	checksum.Write([]byte(header))

	header += fmt.Sprintf("%x%c", checksum.Sum(nil), ENDC)

	var b bytes.Buffer
	b.WriteString(header)
	b.Write(gamefile.Bytes())

	err = ioutil.WriteFile(filename, b.Bytes(), 0644)
	if err != nil {
		log.Fatal("save error:", err)
	}
}

func loadGame(filename string) Game {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("file read error:", err)
	}
	bb := bytes.NewBuffer(dat)
	s, err := bb.ReadString(ENDC)
	if err != nil {
		log.Fatal("buffer read error:", err)
	}
	if len(s) == 0 {
		log.Fatal("file corrupt")
	}
	r, err := zlib.NewReader(bb)
	if err != nil {
		log.Fatal("zlib error:", err)
	}

	dec := gob.NewDecoder(r) // Will read from gamefile.
	if err != nil {
		log.Fatal("decode error 1:", err)
	}

	game := Game{}
	err = dec.Decode(&game)
	if err != nil {
		log.Fatal("decode error 2:", err)
	}
	return game
}

type options struct {
	compile    bool
	run        bool
	save       bool
	file       string
	sourcefile string
	gamefile   string
}

func parseArgs() options {
	usage := `Gork! - Experimental text adventure interpreter

Usage:
  gork <file>
  gork -r <sourcefile>
  gork -c <sourcefile> [-o <gamefile>]
  gork -h | --help | --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  -c --compile  Compile source file.
  -o --outfile  Where to save compiled game.
  -r --run      Run source file directly.`

	args, _ := docopt.ParseArgs(usage, nil, version)

	options := options{}

	options.compile, _ = args.Bool("--compile")
	options.run, _ = args.Bool("--run")
	options.save, _ = args.Bool("--outfile")

	options.file, _ = args.String("<file>")
	options.sourcefile, _ = args.String("<sourcefile>")
	options.gamefile, _ = args.String("<gamefile>")

	return options
}

func runGame(game Game) {
	// game.Pprint(os.Stdout, "")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(game.Title)
	fmt.Println("By", game.By)
	fmt.Println("")
	fmt.Println(game.Description)

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = text[:len(text)-1]

		if strings.Compare("hi", text) == 0 {
			fmt.Println("Hello, Yourself!")
		}
		if strings.Compare("quit", strings.ToLower(text)) == 0 ||
			strings.Compare("goodbye", strings.ToLower(text)) == 0 ||
			strings.Compare("bye", strings.ToLower(text)) == 0 {
			fmt.Println("Goodbye.")
			os.Exit(0)
		}
	}

}

func main() {

	options := parseArgs()

	// run a game file
	if options.file != "" {
		game := loadGame(options.file)
		runGame(game)
		os.Exit(0)
	}

	//compile a source file to a game file
	if options.compile || options.run {

		infile, err := os.Open(options.sourcefile)
		if err != nil {
			log.Fatal("Can't open file:", err)
		}

		game, isErr := parse(infile)
		if isErr {
			os.Exit(1)
		}

		if options.run {
			runGame(game)
			os.Exit(0)
		}

		// otherwise save it
		gamefile := options.gamefile
		if gamefile == "" {
			gamefile = game.Name
		}
		gamefilename := gamefile + ".gork"
		fmt.Printf("Compiling to %s\n", gamefilename)
		saveGame(gamefilename, game)
	}
}
