package parser

import (
	"fmt"

	"github.com/AlexisOMG/compilers-lab7-1/common"
	"github.com/AlexisOMG/compilers-lab7-1/lexer"
)

var (
	Rules = common.Rules{
		common.Expr{
			Kind:  common.NTerm,
			Value: "S",
		}: [][]common.Expr{
			{
				{"AxiomKeyword", common.Term}, {"Nterm", common.Term}, {"NTermKeyword", common.Term}, {"Nterm", common.Term}, {"N", common.NTerm}, {"T", common.NTerm}, {"R", common.NTerm},
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "N",
		}: [][]common.Expr{
			{
				{"Nterm", common.Term}, {"N", common.NTerm},
			},
			{
				common.Epsilon,
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "T",
		}: [][]common.Expr{
			{
				{"TermKeyword", common.Term}, {"Term", common.Term}, {"T1", common.NTerm},
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "T1",
		}: [][]common.Expr{
			{
				{"Term", common.Term}, {"T1", common.NTerm},
			},
			{
				common.Epsilon,
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "R",
		}: [][]common.Expr{
			{
				{"R'", common.NTerm}, {"R1", common.NTerm},
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "R1",
		}: [][]common.Expr{
			{
				{"R'", common.NTerm}, {"R1", common.NTerm},
			},
			{
				common.Epsilon,
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "R'",
		}: [][]common.Expr{
			{
				{"RuleKeyword", common.Term}, {"Nterm", common.Term}, {"Equal", common.Term}, {"V", common.NTerm},
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "V",
		}: [][]common.Expr{
			{
				{"V1", common.NTerm}, {"V2", common.NTerm},
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "V1",
		}: [][]common.Expr{
			{
				{"Term", common.Term}, {"V3", common.NTerm},
			},
			{
				{"Nterm", common.Term}, {"V3", common.NTerm},
			},
			{
				{"EpsKeyword", common.Term},
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "V3",
		}: [][]common.Expr{
			{
				{"Term", common.Term}, {"V3", common.NTerm},
			},
			{
				{"Nterm", common.Term}, {"V3", common.NTerm},
			},
			{
				common.Epsilon,
			},
		},
		common.Expr{
			Kind:  common.NTerm,
			Value: "V2",
		}: [][]common.Expr{
			{
				{"NewLine", common.Term}, {"V", common.NTerm},
			},
			{
				common.Epsilon,
			},
		},
	}

	Terminals = []common.Expr{
		{"AxiomKeyword", common.Term},
		{"NTermKeyword", common.Term},
		{"TermKeyword", common.Term},
		{"RuleKeyword", common.Term},
		{"EpsKeyword", common.Term},
		{"Equal", common.Term},
		{"NewLine", common.Term},
		{"Term", common.Term},
		{"Nterm", common.Term},
	}
)

type Node struct {
	expr     common.Expr
	rule     []common.Expr
	value    string
	Children []*Node
}

func (n *Node) Print() {
	fmt.Print(n.expr.Value, " ")
	if n.expr.Kind == common.NTerm {
		fmt.Print("-> ")
		for _, r := range n.rule {
			fmt.Print(r.Value, " ")
		}
		fmt.Print("\n\tChildren: ")
		for _, child := range n.Children {
			fmt.Print(child.expr.Value, " ")
		}
		fmt.Println()
		for _, child := range n.Children {
			child.Print()
		}
	} else {
		fmt.Println(n.value)
	}
}

type StackItem struct {
	expr   common.Expr
	parent *Node
}

type Stack []StackItem

func Parse(table map[common.Expr]map[common.Expr][][]common.Expr, lex *lexer.Lexer, axiom common.Expr) (*Node, error) {
	var stack Stack
	fakeRoot := Node{
		expr: common.Expr{
			Value: "S'",
			Kind:  common.NTerm,
		},
	}
	stack = append(stack, StackItem{
		expr:   common.Dollar,
		parent: &fakeRoot,
	},
		StackItem{
			expr:   axiom,
			parent: &fakeRoot,
		},
	)

	a := lex.NextToken()
	if a.Kind == lexer.Error {
		return nil, fmt.Errorf("syntax error: %v", a)
	}
	for stack[len(stack)-1].expr != common.Dollar {
		x := stack[len(stack)-1]
		// fmt.Println("STACK: ", stack)
		// fmt.Println("A: ", a, a.Kind.ToString())
		stack = stack[:len(stack)-1]
		if x.expr.Kind == common.Term {
			if x.expr.Value == a.Kind.ToString() {
				x.parent.Children = append(x.parent.Children, &Node{
					expr:  a.ToExpr(),
					value: a.Value,
				})
				a = lex.NextToken()
				if a.Kind == lexer.Error {
					return nil, fmt.Errorf("syntax error: %v", a)
				}
			} else {
				return nil, fmt.Errorf("unexpected %s, expected: %s", a.Kind.ToString(), x.expr.Value)
			}
		} else if exprs := table[x.expr][a.ToExpr()]; exprs[0][0] != common.Error {
			node := Node{
				expr: x.expr,
				rule: exprs[0],
			}
			x.parent.Children = append(x.parent.Children, &node)
			for i := len(exprs[0]) - 1; i >= 0; i-- {
				if exprs[0][i] != common.Epsilon {
					stack = append(stack, StackItem{
						expr:   exprs[0][i],
						parent: &node,
					})
				}
			}
		} else {
			return nil, fmt.Errorf("unexpected %s, expected: %s", a.Kind.ToString(), x.expr.Value)
		}
	}

	// fmt.Println("LAST STACK: ", stack)
	return fakeRoot.Children[0], nil
}
