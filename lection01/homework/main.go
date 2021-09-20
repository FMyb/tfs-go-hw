package main

import "fmt"

type SandglassArgMap map[string]int
type Sandglass func(args SandglassArgMap)

func sandglass(args ...Sandglass) {
	defaultArgs := SandglassArgMap{"size": 15, "color": 0, "char": 'X'}
	for _, arg := range args {
		arg(defaultArgs)
	}
	size := defaultArgs["size"]
	color := defaultArgs["color"]
	char := defaultArgs["char"]
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

func getSandglassSize(size int) Sandglass {
	return func(args SandglassArgMap) {
		args["size"] = size
	}
}

func getSandglassColor(color int) Sandglass {
	return func(args SandglassArgMap) {
		args["color"] = color
	}
}

func getSandglassChar(char int) Sandglass {
	return func(args SandglassArgMap) {
		args["char"] = char
	}
}

func main() {
	sandglass()
	sandglass(getSandglassChar('@'), getSandglassColor(33))
	sandglass(getSandglassSize(16))
}
