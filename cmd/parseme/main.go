package main

import (
	"fmt"
	"os"

	parseme "github.com/fueripe-desu/parseme/pkg"
)

func main() {
	parser := parseme.NewHtmlParser("pkg/parser.go")
	bytes, err := parser.Parse()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	text := string(*bytes)
	fmt.Println(text)
}
