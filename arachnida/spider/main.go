package main

import (
	"fmt"
	"os"
	"net/http"
	"golang.org/x/net/html"
)
func isValidFlag(s string) bool{
	switch s {
	case "-r", "-p", "-l" "-rl", "-rp", "-rlp":
		return true
	}
	return false
}

func main() {
	//Parse (the arguments: -rlp, URL)
	// -rlp, -rlp [N], -rlp [PATH], -rlp [N] [PATH],
	// -r, -rl, -rl [N]
	// -p, -p [PATH]
	if (len(os.Args) < 2){
		fmt.Println("Usage: ./spider URL\nFlags: -rlp")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		rawURL := os.Args[1]
		_, err := http.Get(rawURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		rawURL := os.Args[len(os.Args) -1]
		_, err := http.Get(rawURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		for i := 1; i < len(os.Args) -1; i++ {
			if !isValidFlag(os.Args[i]){
				fmt.Println("Flag's usage: -rlp\nNo repetitions")
				os.Exit(1)
			}
		}
	}
}
