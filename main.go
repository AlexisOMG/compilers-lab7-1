package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlexisOMG/compilers-lab7-1/lexer"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Wrong usage")
	}
	pathToFile := os.Args[1]

	lex, err := lexer.NewLexer(pathToFile)
	if err != nil {
		log.Fatal(err)
	}

	for lex.HasNext() {
		tok := lex.NextToken()
		fmt.Printf("%s %d-%d %s\n", tok.Kind.ToString(), tok.Start, tok.End, tok.Value)
	}
}
