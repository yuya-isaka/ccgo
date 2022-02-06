package codegen

import (
	"fmt"

	h "github.com/yuya-isaka/ccgo/header"
)

func push() {
	fmt.Println("  push %rax")
	h.Depth++
}

func pop(s string) {
	fmt.Printf("  pop %s\n", s)
	h.Depth--
}

func genOffset(s string) int {
	var result int = 0
	var lowerA rune = rune('a')
	for _, v := range s {
		result = int((v - lowerA + 1) * 8)
		break
	}
	return result
}

// %raxにアドレスをセット
func genAddr(node *h.Node) {
	if node.Kind == h.ND_VAR {
		var offset int = genOffset(node.Name)
		fmt.Printf("  lea %d(%%rbp), %%rax\n", -offset)
		return
	}

	panic("not a lvalue")
}

func genExpr(node *h.Node) {
	switch node.Kind {
	case h.ND_NUM:
		fmt.Printf("  mov $%d, %%rax\n", node.Val)
		return
	case h.ND_NEG:
		genExpr(node.Lhs)
		fmt.Println("  neg %rax")
		return
	case h.ND_VAR:
		// メモリにアクセス
		genAddr(node)
		fmt.Println("  mov (%rax), %rax")
		return
	case h.ND_ASSIGN:
		// メモリに格納
		genAddr(node.Lhs)
		push()
		genExpr(node.Rhs)
		pop("%rdi")
		fmt.Println("  mov %rax, (%rdi)")
		return
	}

	genExpr(node.Rhs)
	push()
	genExpr(node.Lhs)
	pop("%rdi")

	switch node.Kind {
	case h.ND_ADD:
		fmt.Println("  add %rdi, %rax")
		return
	case h.ND_SUB:
		fmt.Println("  sub %rdi, %rax")
		return
	case h.ND_MUL:
		fmt.Println("  imul %rdi, %rax")
		return
	case h.ND_DIV:
		fmt.Println("  cqo")
		fmt.Println("  idiv %rdi")
		return
	case h.ND_EQ, h.ND_NE, h.ND_LT, h.ND_LE:
		fmt.Println("  cmp %rdi, %rax")

		if node.Kind == h.ND_EQ {
			fmt.Println("  sete %al")
		} else if node.Kind == h.ND_NE {
			fmt.Println("  setne %al")
		} else if node.Kind == h.ND_LT {
			fmt.Println("  setl %al")
		} else if node.Kind == h.ND_LE {
			fmt.Println("  setle %al")
		}

		fmt.Println("  movzb %al, %rax")
		return
	}

	panic("invalid expression")
}

func genStmt(node *h.Node) {
	if node.Kind == h.ND_EXPR_STMT {
		genExpr(node.Lhs)
		return
	}

	panic("invalid statement")
}

// RBPは現在実行している関数のアドレス（このアドレスの中身は前の関数のアドレス，ポップしたら出てくる）
// RSPはスタックのトップ（こいつに注意）

func prologue() {
	// Prologue (それぞれの関数の冒頭)
	// 関数レコード作成用
	// スタックに今の関数のRBP追加
	fmt.Println("  push %rbp")
	// RSP(今の関数のRBPを指す)の値でRBP更新
	fmt.Println("  mov %rsp, %rbp")
	// RSPから208引いてローカル変数のスタック確保（rspは最後の方をさす）
	fmt.Println("  sub $208, %rsp")
	// 208 == ('z' - 'a' + 1) * 8 ... 64ビット整数,変数26文字のアルファベット用のスタックサイズ
}

func epilogue() {
	// Epilogue
	// 今のRBPの位置にRSPを戻す（確保したローカル変数はここで解放される）
	fmt.Println("  mov %rbp, %rsp")
	// RBPを呼び出し元の関数のRBPに戻す
	fmt.Println("  pop %rbp")
}

func Codegen() {
	fmt.Println("")
	fmt.Println("  .globl main")
	fmt.Println("main:")

	// それぞれの関数の冒頭
	prologue()

	for _, n := range h.Program {
		genStmt(n)
		if h.Depth != 0 {
			panic("wrong")
		}
	}

	epilogue()

	fmt.Println("  ret")
	fmt.Println("")
}
