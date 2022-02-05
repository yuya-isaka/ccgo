package main

import (
	"fmt"
	"os"

	c "github.com/yuya-isaka/ccgo/codegen"
	h "github.com/yuya-isaka/ccgo/header"
	p "github.com/yuya-isaka/ccgo/parse"
	t "github.com/yuya-isaka/ccgo/tokenize"
)

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
		h.Text = append(h.Text, string([]rune{a}[0]))
	}

	h.Tok = t.Tokenize()
	var node *h.Node = p.Expr()

	if !p.EqualKind(h.TK_EOF) {
		h.ErrorToken("extra token")
	}

	fmt.Println("")
	fmt.Println("  .globl main")
	fmt.Println("main:")

	c.GenExpr(node)

	fmt.Println("  ret")
	fmt.Println("")

	if h.Depth != 0 {
		panic("wrong")
	}
}
