package main

import "fmt"

type SandglassArgMap map[string]*int
type Sandglass func(args SandglassArgMap)

func sandglass(args ...Sandglass) {
	size := 15
	color := 0
	var char int = 'X'
	for _, arg := range args {
		arg(SandglassArgMap{"size": &size, "color": &color, "char": &char})
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

func getSandglassSize(size int) Sandglass {
	return func(args SandglassArgMap) {
		sz, ok := args["size"]
		if ok {
			*sz = size
		}
	}
}

func getSandglassColor(color int) Sandglass {
	return func(args SandglassArgMap) {
		clr, ok := args["color"]
		if ok {
			*clr = color
		}
	}
}

func getSandglassChar(char int) Sandglass {
	return func(args SandglassArgMap) {
		chr, ok := args["char"]
		if ok {
			*chr = char
		}
	}
}

func main() {
	sandglass()
	sandglass(getSandglassChar('@'), getSandglassColor(33))
	sandglass(getSandglassSize(16))
}
