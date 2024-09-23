package main
import "fmt";

func main() {

	fmt.Println("Hello World!!!")
	fmt.Println("----------")
	fmt.Printf("Hello %s, you are %d years old.\n", "Alice", 30)
	fmt.Println("----------")
	message := fmt.Sprintf("Hello %s, you are %d years old.", "Alice", 30)
	fmt.Println(message)
}