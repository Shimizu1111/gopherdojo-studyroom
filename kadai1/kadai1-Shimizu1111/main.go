package main

import (
	"flag"
	"fmt"

	img "example.com/image"
)

func main() {
	fmt.Println("start")

	path := flag.String("path", "", "path flag")
	flag.Parse()

	img.ConvertImage(*path)
	fmt.Println("end")
}
