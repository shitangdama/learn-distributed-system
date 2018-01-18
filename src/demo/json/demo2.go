package main

	

import "encoding/json"
import "fmt"
// import "os"
// 
// type Response1 struct {
//     Page   int
//     Fruits []string
// }
// type Response2 struct {
//     Page   int      `json:"page"`
//     Fruits []string `json:"fruits"`
// }
// 可以将map转换成json

type Name struct {
    First, Last string
}

// func main() {
// 	n := Name{First: "Ta", Last: "SY"}

// 	if r, e := json.Marshal(n); e == nil {
// 		fmt.Println("name is ", string(r))
// 	} else {
// 		fmt.Println("err ", e)
// 	}
// }