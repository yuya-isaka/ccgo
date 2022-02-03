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

// 	type Stringer inerface {
//		String() string
//	}
// String()関数を持っているとStringerインターフェースに分類される
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

type Token struct {
	kind TokenKind
	val  int
	loc  int
	str  string
}

func newToken(kind TokenKind) *Token {
	var tok *Token = &Token{} // new(Token)と同じ
	tok.kind = kind
	tok.loc = textNum
	return tok
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
			result = append(result, cur)
			goTok(1)
			continue
		}

		if text[textNum] == "+" || text[textNum] == "-" {
			var cur *Token = newToken(TK_PUNCT)
			cur.str = text[textNum]
			goText(1)
			result = append(result, cur)
			goTok(1)
			continue
		}

		panic("invalid token")
	}

	var cur *Token = newToken(TK_EOF)
	result = append(result, cur)
	goTok(1)

	tokSum = tokNum
	tokNum = 0

	return result
}

func getNumber() int {
	if token[tokNum].kind != TK_NUM {
		panic("Expected a number")
	}
	defer goTok(1)
	return token[tokNum].val
}

func equal(s string) bool {
	return token[tokNum].str == s
}

func skip(s string) {
	if !equal(s) {
		fmt.Fprintf(os.Stdout, "expected %s", s)
		panic("redo")
	}
	defer goTok(1)
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

	fmt.Println("")
	fmt.Println("  .globl main")
	fmt.Println("main:")

	fmt.Printf("  mov $%d, %%rax\n", getNumber())

	for token[tokNum].kind != TK_EOF {
		if equal("+") {
			goTok(1)
			fmt.Printf("  add $%d, %%rax\n", getNumber())
			continue
		}

		skip("-")
		fmt.Printf("  sub $%d, %%rax\n", getNumber())
	}

	fmt.Println("  ret")
	fmt.Println("")
}
