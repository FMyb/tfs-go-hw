package main

import "fmt"

type Sandglass map[string]int

func sandglass(args Sandglass) {
	size := args["size"]
	color := args["color"]
	char, ok := args["char"]
	if !ok {
		char = 'X'
	}
	fmt.Println()
	for i := 0; i < size; i++ {
		fmt.Printf("\033[%dm%c\033[0m", color, char)
	}
	fmt.Println()
	for j := 1; j < size-1; j++ {
		for i := 0; i < size; i++ {
			if i == j || i == size-j-1 {
				fmt.Printf("\u001B[%dm%c\u001B[0m", color, char)
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	for i := 0; i < size; i++ {
		fmt.Printf("\033[%dm%c\033[0m", color, char)
	}
	fmt.Println()
}

func main() {
	sandglass(Sandglass{"size": 20})
	sandglass(Sandglass{"char": '@', "size": 10})
}
