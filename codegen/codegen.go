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

func GenExpr(node *h.Node) {
	switch node.Kind {
	case h.ND_NUM:
		fmt.Printf("  mov $%d, %%rax\n", node.Val)
		return
	case h.ND_NEG:
		GenExpr((node.Lhs))
		fmt.Println("  neg %rax")
		return
	}

	GenExpr(node.Rhs)
	push()
	GenExpr(node.Lhs)
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
