package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/alongubkin/cpl-compiler/pkg/codegen"
	"github.com/alongubkin/cpl-compiler/pkg/parser"
)

func main() {
	// Check args
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "USAGE: ./cpq <input-file>")
		return
	}

	// Make sure the input file ends with .ou
	if path.Ext(os.Args[1]) != ".ou" {
		fmt.Fprintln(os.Stderr, "Input file extension must be .ou")
		return
	}

	// Read code file
	infile := os.Args[1]
	code, err := ioutil.ReadFile(infile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot open input CPL file.")
		return
	}

	// Lex & Parse
	ast, parseErrors := parser.Parse(string(code))
	for _, err := range parseErrors {
		fmt.Fprintf(os.Stderr, "ParseError: %s\n", err.Error())
	}

	// Codegen
	output, codegenErrors := codegen.Codegen(ast)
	for _, err := range codegenErrors {
		fmt.Fprintf(os.Stderr, "CodegenError: %s\n", err.Error())
	}

	// Generate the filename for the output QUAD file
	if len(parseErrors) == 0 && len(codegenErrors) == 0 {
		// Write output to the QUAD file
		outfile := infile[0:len(infile)-3] + ".qud"
		ioutil.WriteFile(outfile, []byte(fixLabels(output)), 0644)
	}
}

func fixLabels(quad string) string {
	labels := 0
	for i, line := range strings.Split(quad, "\n") {
		if strings.HasSuffix(line, ":") {
			label := line[:len(line)-1]
			// Delete label line
			quad = strings.ReplaceAll(quad, line+"\n", "")

			// Replace all label references with the correct line number
			quad = strings.ReplaceAll(quad, label, strconv.Itoa(i-labels+1))
			labels++
		}
	}

	return quad
}
