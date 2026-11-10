package main

import (
	"fmt"
	"os"
	"net/http"
	"golang.org/x/net/html"
	"strconv"
)

func parseDepth(val string) (int, error){
	depth, err := strconv.Atoi(val)
	if err != nil{
		return 0, fmt.Errorf("depth must be a number, got %q", val)
	}
	if depth <= 0{
		return 0, fmt.Errorf("depth must be greater than 0, got %d", depth)
	}
	return depth, nil
}

func isValidFlag(args []string, int leng) bool {
	for  i := 1; i < leng; i++{
		switch args[i] {
		case "-r":
			if (args[i+1]  == "-l"){
				depth, err := parseDepth(args[i+2])
				if err != nil
					return false
				//len(os.Args) == 5 (else error)
				//save depth
			}
			else if (args[i +1]  == "-p"){
				//len(os.Args) == 5 (else error)
				//save the path

			}
			else{
				//len(os.Args) == 4 (else error)
				//recursive default = 5
			}
		case "-p":
			//len(os.Args) == 4 (else error)
			//save the path

		case "-rlp":
			//len(os.Args) == 5 (else error)
			depth, err := parseDepth(args[i+2])
				if err != nil
					return false
			//save depth
			//save the path
		}
	}
}

func main() {
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
		isValidFlag(os.Args, len(os.Args) -1)
		}
}
