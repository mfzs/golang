package main

import (
	"fmt";
	"os"
)

func main(){
	args := os.Args
	fmt.Println("Hello World")
	fmt.Println(args[1])
}
