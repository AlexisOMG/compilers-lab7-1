package main

import (
	"log"
	"os"

	"github.com/AlexisOMG/compilers-lab7-1/common"
	"github.com/AlexisOMG/compilers-lab7-1/lexer"
	"github.com/AlexisOMG/compilers-lab7-1/parser"
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

	rules := parser.Rules

	first := common.First(rules)

	// for l, r := range first {
	// 	fmt.Printf("FIRST(%s) = { ", l.Value)
	// 	for e := range r {
	// 		fmt.Print(e.Value, " ")
	// 	}
	// 	fmt.Println("}")
	// }

	follow := common.Follow(rules, common.Expr{
		Kind:  common.NTerm,
		Value: "S",
	}, first)

	// for l, r := range follow {
	// 	fmt.Printf("FOLLOW(%s) = { ", l.Value)
	// 	for e := range r {
	// 		fmt.Print(e.Value, " ")
	// 	}
	// 	fmt.Println("}")
	// }

	table := common.BuildTable(rules, first, follow, parser.Terminals)

	// fmt.Println()

	// for r, rls := range table {
	// 	fmt.Printf("FOR %v", r)
	// 	for t, exprs := range rls {
	// 		fmt.Printf(" by %v - %v", t, exprs)
	// 	}
	// 	fmt.Println()
	// }

	// fmt.Println()

	answ, err := parser.Parse(table, lex, common.Expr{
		Kind:  common.NTerm,
		Value: "S",
	})
	if err != nil {
		log.Fatal(err)
	}

	answ.Print()
}
