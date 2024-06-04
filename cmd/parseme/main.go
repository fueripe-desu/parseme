package main

import (
	"fmt"
	"os"

	parseme "github.com/fueripe-desu/parseme/pkg"
)

func main() {
	parser := parseme.NewHtmlParser("cmd/parseme/example.html")
	bytes, err := parser.Parse()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	text := *bytes
	fmt.Println(text)
}
