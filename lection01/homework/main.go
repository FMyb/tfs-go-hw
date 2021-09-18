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

func sandglassArg(name string, res int) Sandglass {
	return func() (string, int) {
		return name, res
	}
}

func main() {
	sandglass()
	sandglass(sandglassArg("char", '@'), sandglassArg("color", 33))
	sandglass(sandglassArg("size", 33))
}
