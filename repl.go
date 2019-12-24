package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/despire/interpreter/lexer"
	"github.com/despire/interpreter/token"
)

const (
	prompt = ">>> "
)

// Start reads the input from reader, processes it
// and writes the output to writer.
func Start(reader io.Reader, writer io.Writer) {
	sc := bufio.NewScanner(reader)

	for sc.Scan() {
		if _, err := fmt.Fprintf(writer, prompt); err != nil {
			fmt.Printf("failed to write to output: %+v, reason: %v\n", writer, err)
			continue
		}

		l := lexer.New(sc.Text())
		for tok := l.NextToken(); tok.Typ != token.EOF; tok = l.NextToken() {
			if _, err := fmt.Fprintf(writer, "%+v\n", tok); err != nil {
				fmt.Printf("failed to write to output: %+v, reason: %v\n", writer, err)
			}
		}
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("failed to read from input: %+v\n", reader)
	}
}
