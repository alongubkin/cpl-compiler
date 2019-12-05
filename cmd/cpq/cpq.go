package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/alongubkin/cpl-compiler/pkg/codegen"
	"github.com/alongubkin/cpl-compiler/pkg/parser"
)

func main() {
	// Check args
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "USAGE: ./cpq <input-file>")
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
	ext := path.Ext(os.Args[1])
	outfile := infile[0:len(infile)-len(ext)] + ".qud"

	// Write output to the QUAD file
	ioutil.WriteFile(outfile, []byte(output), 0644)
}
