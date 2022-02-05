package main

import "fmt"

func test() {
	a := []rune("a")
	var lowerA rune = rune('a')
	for _, j := range a {
		if lowerA == j {
			fmt.Println("Yes")
		} else {
			fmt.Println("No")
		}
	}

}
