//gork.go

package main

import (
	"bytes"
	"fmt"
	"os"
    "encoding/gob"
    "log"
	"io/ioutil"
	"compress/zlib"
	"crypto/sha256"
	"time"
)

const fileVersion = "0,0"
const ENDH = '\x0a'
const ENDC = '\x03'

func init(){
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
	for _,s := range in {
		switch n {
		case 0:
			c = "?"
		case 1:
			c = "%"
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
		n +=1 
		if n == 8 {
			n=0
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
		"GORK|%s|%s|%s|%s|%s%c",
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
	log.Fatal("Buffer read error:", err)
	}
	fmt.Println(s)
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

func main() {
	game, err := parse(os.Stdin)

	if !err {
		fmt.Println(game)

		fmt.Println()
		game.Pprint(os.Stdout,"")
		fmt.Println()

		saveGame("/tmp/game.gork",game)

		//
		
		g := loadGame("/tmp/game.gork")

		fmt.Println()
		fmt.Println(g)
		fmt.Println()

		g.Pprint(os.Stdout,"")

	}
}
