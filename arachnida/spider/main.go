
package main
 
import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	// You'll need these once you write the crawler (Go refuses unused imports,
	// so add them back only when you use them):
	// "net/http"
	// "golang.org/x/net/html"
)
 
// Config holds everything the parser found. The crawler only reads this,
// so it never has to look at os.Args again.
type Config struct {
	Recursive bool
	Depth     int    // only meaningful when Recursive is true
	Path      string // where images get saved
	URL       string
}
 
func parseDepth(val string) (int, error) {
	depth, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("depth must be a number, got %q", val)
	}
	if depth <= 0 {
		return 0, fmt.Errorf("depth must be greater than 0, got %d", depth)
	}
	return depth, nil
}
 
// parseArgs keeps your idea: walk the args with a switch on each flag.
// The difference is that each case "eats" the values it needs (i++),
// so the flags work in any order and combination, e.g.
//   ./spider -r URL
//   ./spider -r -l 3 URL
//   ./spider -r -p ./imgs URL
//   ./spider -r -l 3 -p ./imgs URL   (your fixed len() checks couldn't handle this one)
//   ./spider -p ./imgs -r URL
//   ./spider -rlp 3 ./imgs URL
//   ./spider -p ./imgs URL           (no -r: only the given page)
func parseArgs(args []string) (Config, error) {
	cfg := Config{Depth: 5, Path: "./data/"} // defaults from the subject
	depthSet := false
 
	if len(args) < 2 {
		return cfg, fmt.Errorf("missing URL")
	}
 
	last := len(args) - 1 
 
	for i := 1; i < last; i++ {
		switch args[i] {
		case "-r":
			cfg.Recursive = true
 
		case "-l":
			if i+1 >= last { 
				return cfg, fmt.Errorf("-l needs a depth value")
			}
			depth, err := parseDepth(args[i+1])
			if err != nil {
				return cfg, err
			}
			cfg.Depth = depth
			depthSet = true
			i++
 
		case "-p":
			if i+1 >= last {
				return cfg, fmt.Errorf("-p needs a path")
			}
			cfg.Path = args[i+1]
			i++
 
		case "-rlp":
			if i+2 >= last {
				return cfg, fmt.Errorf("-rlp needs a depth and a path: -rlp N PATH URL")
			}
			depth, err := parseDepth(args[i+1])
			if err != nil {
				return cfg, err
			}
			cfg.Recursive = true
			cfg.Depth = depth
			depthSet = true
			cfg.Path = args[i+2]
			i += 2
 
		default:
			return cfg, fmt.Errorf("unknown option %q", args[i])
		}
	}
 
	if depthSet && !cfg.Recursive {
		return cfg, fmt.Errorf("-l can only be used with -r")
	}
 
	raw := args[last]
	if strings.HasPrefix(raw, "-") {
		return cfg, fmt.Errorf("missing URL (last argument is %q)", raw)
	}
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return cfg, fmt.Errorf("invalid URL %q (use http:// or https://)", raw)
	}
	cfg.URL = u.String()
 
	return cfg, nil
}
 
func main() {
	cfg, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Usage: ./spider [-r] [-l N] [-p PATH] URL")
		os.Exit(1)
	}
 
	// Temporary: check the parser before writing anything else.
	fmt.Printf("%+v\n", cfg)
}