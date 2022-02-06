package header

import "fmt"

var TextNum int = 0
var Text []string = make([]string, 0)

var TokSum int = 0
var TokNum int = 0
var Tok []*Token = make([]*Token, 0)

var Depth int = 0

var Program []*Node = make([]*Node, 0)

type TokenKind int

const (
	TK_PUNCT TokenKind = iota
	TK_NUM
	TK_EOF
	TK_IDENT // 変数
)

func (t TokenKind) String() string {
	switch t {
	case TK_PUNCT:
		return "TK_PUNCT"
	case TK_NUM:
		return "TK_NUM"
	case TK_EOF:
		return "TK_EOF"
	case TK_IDENT:
		return "TK_IDENT"
	default:
		return "Unknown"
	}
}

type Token struct {
	Kind TokenKind
	Val  int
	Loc  int
	Str  string
}

func (t *Token) String() string {
	return t.Str
}

type NodeKind int

const (
	ND_ADD       NodeKind = iota // +
	ND_SUB                       // -
	ND_MUL                       // *
	ND_DIV                       // /
	ND_NUM                       // Integer
	ND_NEG                       // -
	ND_EQ                        // ==
	ND_NE                        // !=
	ND_LT                        // <
	ND_LE                        // <=
	ND_EXPR_STMT                 // ;
	ND_ASSIGN                    // =
	ND_VAR                       // Variable
)

// 	type Stringer inerface {
//		String() string
//	}
// String()関数を持っているとStringerインターフェースに分類される
func (t NodeKind) String() string {
	switch t {
	case ND_ADD:
		return "ND_ADD"
	case ND_SUB:
		return "ND_SUB"
	case ND_MUL:
		return "ND_MUL"
	case ND_DIV:
		return "ND_DIV"
	case ND_NUM:
		return "ND_NUM"
	default:
		return "Unknown"
	}
}

type Node struct {
	Kind NodeKind
	Lhs  *Node
	Rhs  *Node
	Val  int    // if Kind == ND_NUM
	Name string // if Kind == ND_VAR
}

func ErrorToken(expect string) {
	fmt.Println()
	for i, v := range Tok {
		fmt.Printf("%dth: %s\n", i, v)
	}
	fmt.Println()
	panic(fmt.Sprintf("%dth token, invalid token, expected %s\n", TokNum, expect))
}

func GoTok(i int) {
	TokNum += i
}

// func getNumber() int {
// 	if Tok[TokNum].Kind != TK_NUM {
// 		ErrorToken("number")
// 	}
// 	defer GoTok(1)
// 	return Tok[TokNum].Val
// }
