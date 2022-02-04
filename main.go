package main

import (
	"fmt"
	"os"
	"strconv"
)

var textNum int = 0
var text []string = make([]string, 0)

var tokSum int = 0
var tokNum int = 0
var token []*Token = make([]*Token, 0)

func goTok(i int) {
	tokNum += i
}

func goText(i int) {
	textNum += i
}

func checkNum() int {
	var result int = 0
	for _, t := range text[textNum:] {
		if v, err := strconv.Atoi(t); err == nil {
			result = result*10 + v
			goText(1)
		} else {
			break
		}
	}
	return result
}

type TokenKind int

const (
	TK_PUNCT TokenKind = iota
	TK_NUM
	TK_EOF
)

func (t TokenKind) String() string {
	switch t {
	case TK_PUNCT:
		return "TK_PUNCT"
	case TK_NUM:
		return "TK_NUM"
	case TK_EOF:
		return "TK_EOF"
	default:
		return "Unknown"
	}
}

type NodeKind int

const (
	ND_ADD NodeKind = iota // +
	ND_SUB                 // -
	ND_MUL                 // *
	ND_DIV                 // /
	ND_NUM                 // Integer
	ND_NEG                 // -
	ND_EQ                  // ==
	ND_NE                  // !=
	ND_LT                  // <
	ND_LE                  // <=
)

type Node struct {
	kind NodeKind
	lhs  *Node
	rhs  *Node
	val  int
}

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

type Token struct {
	kind TokenKind
	val  int
	loc  int
	str  string
}

func (t *Token) String() string {
	return fmt.Sprintf("%s", t.str)
}

func newToken(kind TokenKind) *Token {
	var tok *Token = &Token{} // new(Token)と同じ
	tok.kind = kind
	tok.loc = textNum
	return tok
}

func errorText(s string) {
	fmt.Println()
	for i, v := range text {
		fmt.Printf("%dth: %s\n", i, v)
	}
	fmt.Println()
	panic(fmt.Sprintf("%dth text, %s\n", textNum, s))
}

func ispunctLast(s1 string) int {
	if s1 == "+" || s1 == "-" || s1 == "*" || s1 == "/" || s1 == "(" || s1 == ")" || s1 == "<" || s1 == ">" || s1 == "=" {
		return 1
	} else {
		return 0
	}
}

func ispunct(s1, s2 string) int {
	if (s1+s2) == "==" || (s1+s2) == "!=" || (s1+s2) == "<=" || (s1+s2) == ">=" {
		return 2
	} else {
		return ispunctLast(s1)
	}
}

func tokenize() []*Token {
	var result []*Token = make([]*Token, 0)

	for len(text) > textNum {
		if text[textNum] == " " {
			goText(1)
			continue
		}

		if _, err := strconv.Atoi(text[textNum]); err == nil {
			var cur *Token = newToken(TK_NUM)
			var tmp int = checkNum()
			cur.val = tmp
			cur.str = strconv.Itoa(tmp)
			result = append(result, cur)
			goTok(1)
			continue
		}

		var flag int = 0
		if len(text)-1 > textNum {
			flag = ispunct(text[textNum], text[textNum+1])
		} else {
			flag = ispunctLast(text[textNum])
		}
		if flag == 1 || flag == 2 {
			var cur *Token = newToken(TK_PUNCT)
			cur.str = text[textNum]
			if flag == 2 {
				cur.str += text[textNum+1]
			}
			goText(flag)
			result = append(result, cur)
			goTok(1)
			continue
		}

		errorText("invalid text")
	}

	var cur *Token = newToken(TK_EOF)
	cur.str = "EOF"
	result = append(result, cur)
	goTok(1)

	tokSum = tokNum
	tokNum = 0

	return result
}

func errorToken(expect string) {
	fmt.Println()
	for i, v := range token {
		fmt.Printf("%dth: %s\n", i, v)
	}
	fmt.Println()
	panic(fmt.Sprintf("%dth token, invalid token, expected %s\n", tokNum, expect))
}

func getNumber() int {
	if token[tokNum].kind != TK_NUM {
		errorToken("number")
	}
	defer goTok(1)
	return token[tokNum].val
}

func equalStr(s string) bool {
	return token[tokNum].str == s
}

func equalKind(t TokenKind) bool {
	return token[tokNum].kind == t
}

func equalStrGo(s string) bool {
	if token[tokNum].str == s {
		goTok(1)
		return true
	} else {
		return false
	}
}

func hopeStrGo(s string) {
	if !equalStr(s) {
		errorToken(s)
	}
	defer goTok(1)
}

func expr() *Node {
	return equality()
}

func equality() *Node {
	var node *Node = relational()

	for {
		if equalStrGo("==") {
			node = newBinary(ND_EQ, node, relational())
			continue
		}

		if equalStrGo("!=") {
			node = newBinary(ND_NE, node, relational())
			continue
		}

		return node
	}
}

func relational() *Node {
	var node *Node = add()

	for {
		if equalStrGo("<") {
			node = newBinary(ND_LT, node, add())
			continue
		}

		if equalStrGo("<=") {
			node = newBinary(ND_LE, node, add())
			continue
		}

		if equalStrGo(">") {
			node = newBinary(ND_LT, add(), node)
			continue
		}

		if equalStrGo(">=") {
			node = newBinary(ND_LE, add(), node)
			continue
		}

		return node
	}
}

func add() *Node {
	var node *Node = mul()

	for {
		if equalStrGo("+") {
			node = newBinary(ND_ADD, node, mul())
			continue
		}

		if equalStrGo("-") {
			node = newBinary(ND_SUB, node, mul())
			continue
		}

		return node
	}
}

func mul() *Node {
	var node *Node = unary()

	for {
		if equalStrGo("*") {
			node = newBinary(ND_MUL, node, unary())
			continue
		}

		if equalStrGo("/") {
			node = newBinary(ND_DIV, node, unary())
			continue
		}

		return node
	}
}

func unary() *Node {
	if equalStrGo("+") {
		return unary()
	}

	if equalStrGo("-") {
		return newUnary(ND_NEG, unary())
	}

	return primary()
}

func primary() *Node {
	if equalStrGo("(") {
		var node *Node = expr()
		hopeStrGo(")")
		return node
	}

	if equalKind(TK_NUM) {
		var node *Node = newNum(token[tokNum].val)
		defer goTok(1)
		return node
	}

	errorToken("expected an expression")
	return newNode(ND_NUM) // ここは実行されないはず
}

func newNode(kind NodeKind) *Node {
	var node *Node = &Node{}
	node.kind = kind
	return node
}

func newNum(val int) *Node {
	var node *Node = newNode(ND_NUM)
	node.val = val
	return node
}

func newBinary(kind NodeKind, lhs *Node, rhs *Node) *Node {
	var node *Node = newNode(kind)
	node.lhs = lhs
	node.rhs = rhs
	return node
}

func newUnary(kind NodeKind, expr *Node) *Node {
	var node *Node = newNode(kind)
	node.lhs = expr
	return node
}

var depth int = 0

func push() {
	fmt.Println("  push %rax")
	depth++
}

func pop(s string) {
	fmt.Printf("  pop %s\n", s)
	depth--
}

func genExpr(node *Node) {
	switch node.kind {
	case ND_NUM:
		fmt.Printf("  mov $%d, %%rax\n", node.val)
		return
	case ND_NEG:
		genExpr((node.lhs))
		fmt.Println("  neg %rax")
		return
	}

	genExpr(node.rhs)
	push()
	genExpr(node.lhs)
	pop("%rdi")

	switch node.kind {
	case ND_ADD:
		fmt.Println("  add %rdi, %rax")
		return
	case ND_SUB:
		fmt.Println("  sub %rdi, %rax")
		return
	case ND_MUL:
		fmt.Println("  imul %rdi, %rax")
		return
	case ND_DIV:
		fmt.Println("  cqo")
		fmt.Println("  idiv %rdi")
		return
	case ND_EQ, ND_NE, ND_LT, ND_LE:
		fmt.Println("  cmp %rdi, %rax")

		if node.kind == ND_EQ {
			fmt.Println("  sete %al")
		} else if node.kind == ND_NE {
			fmt.Println("  setne %al")
		} else if node.kind == ND_LT {
			fmt.Println("  setl %al")
		} else if node.kind == ND_LE {
			fmt.Println("  setle %al")
		}

		fmt.Println("  movzb %al, %rax")
		return
	}

	panic("invalid expression")
}

func main() {
	// a := "あ"
	// "あ" -> 3042 ... 文字コードの規格Unicodeで決められたcode point
	// 3042 -> E3 81 82 ... UTF-8符号化方式でcode pointを1byteから4byteの可変長のバイトデータでに置換
	// a[0] は E3 ... 文字列へのインデックスアクセスはバイト単位でのアクセス
	// 文字列の長さを知りたいときは，len(a)ではなくlen([]rune{a})
	// argLength := len(os.Args[1:])
	var argLength int = len(os.Args[1:])
	if argLength != 1 {
		fmt.Printf("Arg length is %d\n", argLength)
		panic("invalid argment number")
	}

	for _, a := range os.Args[1] {
		text = append(text, string([]rune{a}[0]))
	}

	token = tokenize()
	var node *Node = expr()

	if token[tokNum].kind != TK_EOF {
		errorToken("extra token")
	}

	fmt.Println("")
	fmt.Println("  .globl main")
	fmt.Println("main:")

	genExpr(node)

	fmt.Println("  ret")
	fmt.Println("")

	if depth != 0 {
		panic("wrong")
	}
}
