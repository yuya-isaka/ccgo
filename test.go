package main

import "fmt"

func test() {
	a := "b"
	var lowerA rune = rune('a')
	for _, v := range a {
		fmt.Println((v - lowerA + 1) * 8)
	}
}
