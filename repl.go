package main

import (
	"bufio"
	"fmt"
	"github.com/despire/interpreter/eval"
	"github.com/despire/interpreter/objects"
	"github.com/despire/interpreter/parser"
	"io"

	"github.com/despire/interpreter/lexer"
)

const (
	prompt = ">>> "
)

// Start reads the input from reader, processes it
// and writes the output to writer.
func Start(reader io.Reader, writer io.Writer) {
	sc := bufio.NewScanner(reader)
	env := objects.NewEnvironment()

	for sc.Scan() {
		if _, err := fmt.Fprintf(writer, prompt); err != nil {
			fmt.Printf("failed to write to output: %+v, reason: %v\n", writer, err)
			continue
		}

		l := lexer.New(sc.Text())
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(writer, p.Errors())
			continue
		}

		e := eval.Eval(program, env)
		if e != nil {
			io.WriteString(writer, e.Inspect())
			io.WriteString(writer, "\n")
		}
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("failed to read from input: %+v\n", reader)
	}
}

func printParseErrors(writer io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(writer, "\t"+err+"\n")
	}
}
