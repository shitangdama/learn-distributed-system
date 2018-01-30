package main



import (
	"fmt"
)


type Raft struct {


}


func main() {

	fmt.Println("开始test")

	rf := &Raft{}
	go rf.change()
}

func (rf *Raft) change() {

}