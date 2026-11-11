package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	//"github.com/bep/imagemeta"
)

var allowedExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
}

func parseArgs(args []string) error {
	for i := 1; i < len(args); i++ {
		info, err := os.Stat(args[i])
		if err != nil {
			return fmt.Errorf("non valid path")
		}
		if info.IsDir() {
			return fmt.Errorf("no file")
		}
		ext := strings.ToLower(filepath.Ext(args[i]))
		if !allowedExt[ext] {
			return fmt.Errorf("%s: unsupported extension %q", args[i], ext)
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: missing files")
		fmt.Fprintln(os.Stderr, "Usage: ./scorpion FILE1 [FILE2 ...]")
		os.Exit(1)
	}
	err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Usage: ./scorpion FILE1 [FILE2 ...]")
		os.Exit(1)
	}

}

// f, err := os.Open("image.png")
// if err != nil {
// 	panic(err)
// }
// defer f.Close()
