package main

import (
	"fmt"
	"os"
)

func test() {
	for _, a := range os.Args[1] {
		text = append(text, string([]rune{a}[0]))
	}
	for _, a := range text {
		if a == " " {
			fmt.Println("No!!!")
		}
	}
}
