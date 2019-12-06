# CPL Compiler

This project compiles CPL, which is a very small subset of C, to an assembly language called Quad. The compiler is written in Go.

This is my project for the [Compilation Course](https://www.openu.ac.il/courses/20364.htm) in the Open University of Israel.

## Usage

To compile a CPL file, simply run the following command:

    cpq myfile.ou

The program will create a `.qud` output file in the same directory.

## Building and Testing

### Requirements

The only requirement for building the project is [Go](https://golang.org/).

### Compilation

The following script will cross-compile the project for all platforms:

    ./cross-compile.sh

### Running Tests

    go test ./pkg/lexer
    go test ./pkg/parser
    go test ./pkg/codegen