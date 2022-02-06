package parse

import (
	h "github.com/yuya-isaka/ccgo/header"
)

// プログラム = 文*
// 文 = 式;
// 式 = 等値判定 (= 等値判定)?
// 等値判定 = 大小比較 (== 大小比較 | != 大小比較)*
// 大小比較 = 加減 ("<" 加減 | "<=" 加減 | ">" 加減 | ">=" 加減)*
// 加減 = 乗除 ("+" 乗除 | "/" 乗除)*
// 乗除 = 単項 ("*" 単項 | "/" 単項)*
// 単項 = ("+" | "-")? 終端
// 終端 = 整数 | 変数 | "(" 式 ")"

// 評価の優先順位は下から上

// 式 ＝ 等値判定 (= 等値判定)?
// 代入式になる場合は，式の部分は左辺値でなければならない
// 左辺値とは基本的にメモリのアドレスを指定する式のみ（レジスタに存在していたり，固定のオフセットでアクセスできなかったらダメ）

func equalStr(s string) bool {
	return h.Tok[h.TokNum].Str == s
}

func equalKind(tk h.TokenKind) bool {
	return h.Tok[h.TokNum].Kind == tk
}

func equalStrGo(s string) bool {
	if h.Tok[h.TokNum].Str == s {
		h.GoTok(1)
		return true
	} else {
		return false
	}
}

func hopeStrGo(s string) {
	if !equalStr(s) {
		h.ErrorToken(s)
	}
	defer h.GoTok(1)
}

func newNode(kind h.NodeKind) *h.Node {
	var node *h.Node = &h.Node{}
	node.Kind = kind
	return node
}

func newNum(val int) *h.Node {
	var node *h.Node = newNode(h.ND_NUM)
	node.Val = val
	return node
}

func newIdent(name string) *h.Node {
	var node *h.Node = newNode(h.ND_VAR)
	node.Name = name
	return node
}

func newBinary(kind h.NodeKind, lhs *h.Node, rhs *h.Node) *h.Node {
	var node *h.Node = newNode(kind)
	node.Lhs = lhs
	node.Rhs = rhs
	return node
}

func newUnary(kind h.NodeKind, expr *h.Node) *h.Node {
	var node *h.Node = newNode(kind)
	node.Lhs = expr
	return node
}

func Parse() []*h.Node {
	var result []*h.Node = make([]*h.Node, 0)
	for !equalKind(h.TK_EOF) {
		result = append(result, stmt())
	}
	return result
}

func stmt() *h.Node {
	var node *h.Node = newUnary(h.ND_EXPR_STMT, assign())
	hopeStrGo(";")
	return node
}

func assign() *h.Node {
	var node *h.Node = expr()

	if equalStrGo("=") {
		// assignの後にassignが来てもいい
		node = newBinary(h.ND_ASSIGN, node, assign())
	}

	return node
}

func expr() *h.Node {
	var node *h.Node = relational()

	for {
		if equalStrGo("==") {
			node = newBinary(h.ND_EQ, node, relational())
			continue
		}

		if equalStrGo("!=") {
			node = newBinary(h.ND_NE, node, relational())
			continue
		}

		return node
	}
}

func relational() *h.Node {
	var node *h.Node = add()

	for {
		if equalStrGo("<") {
			node = newBinary(h.ND_LT, node, add())
			continue
		}

		if equalStrGo("<=") {
			node = newBinary(h.ND_LE, node, add())
			continue
		}

		if equalStrGo(">") {
			node = newBinary(h.ND_LT, add(), node)
			continue
		}

		if equalStrGo(">=") {
			node = newBinary(h.ND_LE, add(), node)
			continue
		}

		return node
	}
}

func add() *h.Node {
	var node *h.Node = mul()

	for {
		if equalStrGo("+") {
			node = newBinary(h.ND_ADD, node, mul())
			continue
		}

		if equalStrGo("-") {
			node = newBinary(h.ND_SUB, node, mul())
			continue
		}

		return node
	}
}

func mul() *h.Node {
	var node *h.Node = unary()

	for {
		if equalStrGo("*") {
			node = newBinary(h.ND_MUL, node, unary())
			continue
		}

		if equalStrGo("/") {
			node = newBinary(h.ND_DIV, node, unary())
			continue
		}

		return node
	}
}

func unary() *h.Node {
	if equalStrGo("+") {
		return unary()
	}

	if equalStrGo("-") {
		return newUnary(h.ND_NEG, unary())
	}

	return primary()
}

func primary() *h.Node {
	if equalStrGo("(") {
		var node *h.Node = expr()
		hopeStrGo(")")
		return node
	}

	if equalKind(h.TK_NUM) {
		var node *h.Node = newNum(h.Tok[h.TokNum].Val)
		defer h.GoTok(1)
		return node
	}

	if equalKind(h.TK_IDENT) {
		var node *h.Node = newIdent(h.Tok[h.TokNum].Str)
		defer h.GoTok(1)
		return node
	}

	h.ErrorToken("expected an expression")
	return newNode(h.ND_NUM) // ここは実行されないはず
}
