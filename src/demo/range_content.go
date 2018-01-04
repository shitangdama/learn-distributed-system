package main

import (
	"fmt"
)



func main(){
	contents := "as fsf af: afasdf"
	for i, c := range contents {
		fmt.Printf("%d",i)
		fmt.Printf(string(c))
	}
}