package main

import "fmt"

type Sandglass func() (string, int)

func sandglass(args ...Sandglass) {
	size := 15
	color := 0
	var char int = 'X'
	for _, arg := range args {
		name, res := arg()
		switch name {
		case "size":
			size = res
		case "char":
			char = res
		case "color":
			color = res
		}
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

func getSandglassChar(char int) Sandglass {
	return func() (string, int) {
		return "char", char
	}
}

func getSandglassSize(size int) Sandglass {
	return func() (string, int) {
		return "size", size
	}
}

func getSandglassColor(color int) Sandglass {
	return func() (string, int) {
		return "color", color
	}
}

func main() {
	sandglass()
	sandglass(getSandglassChar('@'), getSandglassColor(33))
}
