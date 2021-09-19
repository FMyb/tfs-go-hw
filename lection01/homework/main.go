package main

import "fmt"

type Sandglass func(*int, *int, *int)

func sandglass(args ...Sandglass) {
	size := 15
	color := 0
	var char int = 'X'
	for _, arg := range args {
		arg(&size, &color, &char)
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
	return func(sizeArg *int, colorArg *int, charArg *int) {
		*charArg = char
	}
}

func getSandglassSize(size int) Sandglass {
	return func(sizeArg *int, colorArg *int, charArg *int) {
		*sizeArg = size
	}
}

func getSandglassColor(color int) Sandglass {
	return func(sizeArg *int, colorArg *int, charArg *int) {
		*colorArg = color
	}
}

func main() {
	sandglass()
	sandglass(getSandglassChar('@'), getSandglassColor(33))
	sandglass(getSandglassSize(16))
}
