package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
)

var (
	inputFile  = flag.String("input", "", "Input file containing SubHandler definition")
	outputFile = flag.String("output", "", "Output file to write handler function to")
)

func init() {
	flag.Parse()
}

func main() {
	if len(*inputFile) == 0 || len(*outputFile) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Generating handler definitions, %s -> %s\n", *inputFile, *outputFile)

	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, *inputFile, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Parse input file: %v\n", err)
	}

	// Find the SubHandler struct.
	subHandler, ok := getSubHandler(node)
	if !ok {
		log.Fatalf("Failed to find SubHandler struct in input file\n")
	}

	// Analyze fields of SubHandler struct to build cases
	cases := buildHandlerCases(subHandler)

	// Create output file
	output, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Create output file: %v\n", err)
	}
	defer output.Close()

	// Render template to output file
	err = renderOutputFile(cases, output)
	if err != nil {
		log.Fatalf("Render output file: %v\n", err)
	}
}
