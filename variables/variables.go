package main
import "fmt"

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

}