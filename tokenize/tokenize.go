package tokenize

import (
	"fmt"
	"strconv"

	h "github.com/yuya-isaka/ccgo/header"
)

func newToken(kind h.TokenKind) *h.Token {
	var tok *h.Token = &h.Token{} // new(Token)と同じ
	tok.Kind = kind
	tok.Loc = h.TextNum
	return tok
}

func goText(i int) {
	h.TextNum += i
}

func checkNum() int {
	var result int = 0
	for _, t := range h.Text[h.TextNum:] {
		if v, err := strconv.Atoi(t); err == nil {
			result = result*10 + v
			goText(1)
		} else {
			break
		}
	}
	return result
}

func errorText(s string) {
	fmt.Println()
	for i, v := range h.Text {
		fmt.Printf("%dth: %s\n", i, v)
	}
	fmt.Println()
	panic(fmt.Sprintf("%dth text, %s\n", h.TextNum, s))
}

func ispunctLast(s1 string) int {
	if s1 == "+" || s1 == "-" || s1 == "*" || s1 == "/" || s1 == "(" || s1 == ")" || s1 == "<" || s1 == ">" || s1 == "=" || s1 == ";" {
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

func isAlphabet(s string) bool {
	var result bool = false
	var lowerA rune = rune('a')
	var lowerZ rune = rune('z')
	for _, v := range s {
		if lowerA <= v && v <= lowerZ {
			result = true
			break
		} else {
			result = false
			break
		}
	}
	return result
}

func Tokenize() []*h.Token {
	var result []*h.Token = make([]*h.Token, 0)

	for len(h.Text) > h.TextNum {
		if h.Text[h.TextNum] == " " {
			goText(1)
			continue
		}

		if _, err := strconv.Atoi(h.Text[h.TextNum]); err == nil {
			var cur *h.Token = newToken(h.TK_NUM)
			var tmp int = checkNum()
			cur.Val = tmp
			cur.Str = strconv.Itoa(tmp)
			result = append(result, cur)
			h.GoTok(1)
			continue
		}

		var flag int = 0
		if len(h.Text)-1 > h.TextNum {
			flag = ispunct(h.Text[h.TextNum], h.Text[h.TextNum+1])
		} else {
			flag = ispunctLast(h.Text[h.TextNum])
		}
		if flag == 1 || flag == 2 {
			var cur *h.Token = newToken(h.TK_PUNCT)
			cur.Str = h.Text[h.TextNum]
			if flag == 2 {
				cur.Str += h.Text[h.TextNum+1]
			}
			goText(flag)
			result = append(result, cur)
			h.GoTok(1)
			continue
		}

		if isAlphabet(h.Text[h.TextNum]) {
			// ローカル変数として保存
			var cur *h.Token = newToken(h.TK_RESERVED)
			cur.Str = h.Text[h.TextNum]
			goText(1)

			// トークン列更新
			result = append(result, cur)
			h.GoTok(1)

			continue
		}

		errorText("invalid text")
	}

	var cur *h.Token = newToken(h.TK_EOF)
	cur.Str = "EOF"
	result = append(result, cur)
	h.GoTok(1)

	h.TokSum = h.TokNum
	h.TokNum = 0

	return result
}
