package main

import (
	"fmt"
	"os"
	"strconv"
)

var num int = 0
var text []string = make([]string, 0)

func goNum(i int) {
	num += i
}

func checkNum() int {
	var result int = 0
	for _, t := range text[num:] {
		if v, err := strconv.Atoi(t); err == nil {
			result = result*10 + v
			goNum(1)
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
	next *Token
	val  int
	loc  int
	len  int
	str  string
}

func newToken(kind TokenKind) *Token {
	var tok *Token = &Token{} // new(Token)と同じ
	tok.kind = kind
	tok.loc = num
	return tok
}

func tokenize() *Token {
	var head Token = Token{}
	var cur *Token = &head

	for len(text) > num {
		if text[num] == " " {
			goNum(1)
			continue
		}

		if _, err := strconv.Atoi(text[num]); err == nil {
			cur.next = newToken(TK_NUM)
			cur = cur.next
			var tmp int = checkNum()
			cur.val = tmp
			continue
		}

		if text[num] == "+" || text[num] == "-" {
			cur.next = newToken(TK_PUNCT)
			cur = cur.next
			cur.str = text[num]
			goNum(1)
			continue
		}

		panic("invalid token")
	}

	cur.next = newToken(TK_EOF)
	cur = cur.next

	return head.next
}

func getNumber(tok *Token) int {
	if tok.kind != TK_NUM {
		panic("Expected a number")
	}
	return tok.val
}

func equal(tok *Token, s string) bool {
	return tok.str == s
}

func skip(tok *Token, s string) *Token {
	if !equal(tok, s) {
		fmt.Fprintf(os.Stdout, "expected %s", s)
		panic("redo")
	}
	return tok.next
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

	var tok *Token = tokenize()

	fmt.Println("")
	fmt.Println("  .globl main")
	fmt.Println("main:")

	fmt.Printf("  mov $%d, %%rax\n", getNumber(tok))
	tok = tok.next

	for tok.kind != TK_EOF {
		if equal(tok, "+") {
			fmt.Printf("  add $%d, %%rax\n", getNumber(tok.next))
			tok = tok.next.next
			continue
		}

		tok = skip(tok, "-")
		fmt.Printf("  sub $%d, %%rax\n", getNumber(tok))
		tok = tok.next
	}

	fmt.Println("  ret")
	fmt.Println("")
}
