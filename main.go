package loggerhead

import (
	"flag"
)

var (
	workDir     = flag.String("f", ".", "Working directory to be used during runtime")
	program     = flag.String("b", "", "Program to start in host OS binary format for instance python")
	programArgs = flag.String("a", "", "space separated arguement to be pass to the program for instance python filename.py")
	stdout      = flag.String("o", "", "file path to show stdout and stderror result. default to stdout.in home directory")
)

func main() {
	flag.Parse()
}
