package demo

import (
	"fmt",
	"unicode"
)



func main() {
	s := "Hello\n\t世界！"
	for _, r := range s {
	fmt.Printf("%c = %v\n", r, unicode.IsLetter(r))
	} // Hello世界 = true
}