package main

import (
	"fmt"
	"os"
)

func main() {
	for _, s := range os.Args {
		fmt.Println(s)
	}
}
