package main

import (
	"fmt"
	"os"
)

func reverse(s string) string { 
    rns := []rune(s) // convert to rune 
    for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 { 
        // swap the letters of the string, 
        // like first with last and so on. 
        rns[i], rns[j] = rns[j], rns[i] 
    } 
  
    // return the reversed string. 
    return string(rns) 
} 

func main(){
	fmt.Println("Hi")
	a := 5
	b := 7
	var c float64  = 9.2
	var d float64  = 11.1

	fmt.Println("-----------------------------------------")
	fmt.Printf("The sum of %d and %d is: %d\n", a,b,a+b)
	fmt.Printf("The diff of %d and %d is: %d\n", a,b,a-b)
	fmt.Printf("The div of %d and %d is: %d\n", a,b,a/b)
	fmt.Printf("The actual div of %d and %d is: %f\n", a,b, float64(a)/float64(b))
	fmt.Printf("The multiply of %d and %d is: %d\n", a,b,a*b)

	fmt.Println("-----------------------------------------")
	fmt.Printf("The sum of %f and %f is: %f\n", c,d,c+d)
	fmt.Printf("The diff of %f and %f is: %.2f\n", c,d,c-d)
	fmt.Printf("The multiply of %f and %f is: %.2f\n", c,d,c*d)
	fmt.Printf("The div of %f and %f is: %.2f\n", c,d,c/d)

	fmt.Println("-----------------------------------------")
	var e = "My name is"
	var f = "Firoz"
	fmt.Printf(e+" "+f)

	fmt.Println("-----------------------------------------")
	newslice := e[0:7]
	fmt.Println(newslice)
	// newreverse :=[:-1]
	// fmt.Printf(newreverse)
	strRev := reverse(f) 
    fmt.Println(f) 
    fmt.Println(strRev)

	// Calling Args from command lines
	args := os.Args
	fmt.Printf("Showing the command line Args %v",args[1:]) //len(args)
}