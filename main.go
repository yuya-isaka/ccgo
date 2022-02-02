package main

import (
	"fmt"
	"os"
	"strconv"
)

var num int = 0

func checkNum(s []string) int {
	tmp := 0
	for _, t := range s {
		if v, err := strconv.Atoi(t); err == nil {
			tmp = tmp*10 + v
			goNum(1)
		} else {
			break
		}
	}
	return tmp
}

func goNum(i int) {
	num += i
}

func main() {
	var text []string = make([]string, 0)
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

	fmt.Println("")
	fmt.Println("  .globl main")
	fmt.Println("main:")
	tmp := checkNum(text[num:])
	fmt.Printf("  mov $%d, %%rax\n", tmp)

	for len(text) > num {
		if text[num] == "+" {
			goNum(1)
			tmp := checkNum(text[num:])
			fmt.Printf("  add $%d, %%rax\n", tmp)
			continue
		}

		if text[num] == "-" {
			goNum(1)
			tmp := checkNum(text[num:])
			fmt.Printf("  sub $%d, %%rax\n", tmp)
			continue
		}
		panic("unexpected character")
	}

	fmt.Println("  ret")
	fmt.Println("")
}
