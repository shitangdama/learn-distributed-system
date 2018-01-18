package main

import (
	"fmt"
	// "strings" //io 工具包
	"encoding/json"
)

// marshal用于把各个函数转换成json
func main() {
	bolB, _ := json.Marshal(true)  
	fmt.Println(string(bolB))

	intB, _ := json.Marshal(1)  
	fmt.Println(string(intB))

	fltB, _ := json.Marshal(2.34) 
	fmt.Println(string(fltB))

	strB, _ := json.Marshal("gopher") 
	fmt.Println(string(strB))

	slcD := []string{"apple", "peach", "pear"} 
	slcB, _ := json.Marshal(slcD) 
	fmt.Println(string(slcB))

	mapD := map[string]int{"apple": 5, "lettuce": 7} 
	mapB, _ := json.Marshal(mapD) 
	fmt.Println(string(mapB))
}