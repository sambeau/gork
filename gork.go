//gork.go

package main

import (
	"fmt"
	"os"
)

func main() {
	game, err := parse(os.Stdin)

	if !err {
		fmt.Println(game)
	}
}
