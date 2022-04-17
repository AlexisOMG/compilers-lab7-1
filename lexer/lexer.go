package lexer

import (
	"io/ioutil"
	"regexp"
)

var (
	axiomKeywordReg = regexp.MustCompile(`^\$AXIOM`)
	ntermKeywordReg = regexp.MustCompile(`^\$NTERM`)
	termKeywordReg  = regexp.MustCompile(`^\$TERM`)
	ruleKeywordReg  = regexp.MustCompile(`^\$RULE`)
	epsKeywordReg   = regexp.MustCompile(`^\$EPS`)
	ntermReg        = regexp.MustCompile(`^[A-Z][^ \n]*`)
	termReg         = regexp.MustCompile(`^"[^ \nA-Z]+"`)
	equalReg        = regexp.MustCompile(`^=`)
	newLineReg      = regexp.MustCompile(`^\n`)
)

const (
	AxiomKeyword = iota
	NTermKeyword
	TermKeyword
	RuleKeyword
	EpsKeyword
	Term
	Nterm
	Equal
	NewLine
	EOF
	Error
)

type regWithKind struct {
	kind Kind
	reg  *regexp.Regexp
}

type Kind int

func (k Kind) ToString() string {
	switch k {
	case AxiomKeyword:
		return "AxiomKeyword"
	case NTermKeyword:
		return "NTermKeyword"
	case TermKeyword:
		return "TermKeyword"
	case RuleKeyword:
		return "RuleKeyword"
	case EpsKeyword:
		return "EpsKeyword"
	case Term:
		return "Term"
	case Nterm:
		return "Nterm"
	case Equal:
		return "Equal"
	case NewLine:
		return "NewLine"
	case EOF:
		return "EOF"
	case Error:
		return "Error"
	}

	return "unknown kind"
}

type Token struct {
	Kind  Kind
	Value string
	Start int
	End   int
}

type Lexer struct {
	text     string
	regs     []regWithKind
	curIndex int
}

func (l *Lexer) HasNext() bool {
	return len(l.text) > 0
}

func (l *Lexer) NextToken() Token {
	if !l.HasNext() {
		return Token{
			Kind:  EOF,
			Start: l.curIndex + 1,
			End:   l.curIndex + 1,
		}
	}

	if l.text[0] == ' ' || l.text[0] == '\t' {
		l.text = l.text[1:]
		l.curIndex += 1
		return l.NextToken()
	}

	for _, r := range l.regs {
		if loc := r.reg.FindStringIndex(l.text); loc != nil {
			token := Token{
				Kind:  r.kind,
				Value: l.text[loc[0]:loc[1]],
				Start: l.curIndex,
				End:   l.curIndex + loc[1] - loc[0],
			}
			if token.Kind == NewLine {
				token.Value = `\n`
			}
			l.text = l.text[loc[1]:]
			l.curIndex += (loc[1] - loc[0])
			return token
		}
	}

	tok := Token{
		Kind:  Error,
		Start: l.curIndex,
		End:   l.curIndex,
	}

	l.curIndex += 1
	l.text = l.text[1:]

	return tok
}

func NewLexer(pathToFile string) (*Lexer, error) {
	data, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	return &Lexer{
		text:     string(data),
		curIndex: 1,
		regs: []regWithKind{
			{
				reg:  axiomKeywordReg.Copy(),
				kind: AxiomKeyword,
			},
			{
				reg:  ntermKeywordReg.Copy(),
				kind: NTermKeyword,
			},
			{
				reg:  termKeywordReg.Copy(),
				kind: TermKeyword,
			},
			{
				reg:  ruleKeywordReg.Copy(),
				kind: RuleKeyword,
			},
			{
				reg:  epsKeywordReg.Copy(),
				kind: EpsKeyword,
			},
			{
				reg:  ntermReg.Copy(),
				kind: Nterm,
			},
			{
				reg:  termReg.Copy(),
				kind: Term,
			},
			{
				reg:  equalReg.Copy(),
				kind: Equal,
			},
			{
				reg:  newLineReg.Copy(),
				kind: NewLine,
			},
		},
	}, nil
}
