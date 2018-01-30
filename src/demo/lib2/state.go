package main



import (
	"fmt"
	"sync"
	"time"
	"math/rand"
)

const (
	STATE_LEADER = iota
	STATE_CANDIDATE
	STATE_FOLLOWER
)

type Raft struct {
	mu sync.Mutex
	state int
	time *time.Timer
	requestForVoteChan chan bool

}


func main() {

	fmt.Println("开始test")

	rf := &Raft{}
	rf.state = STATE_LEADER
	// receivedHeartbeat bool
	rf.change()
}


func (rf *Raft) change(){
	fmt.Println("zhuangtao1")
	m:=rand.Intn(150)
    m+=150
	// rf.time = time.NewTimer()
	for{
		select{
		case <-rf.requestForVoteChan:
            return

		case <-time.After(time.Duration(m)*time.Millisecond):
			fmt.Println("xintiao")
		}
	}

}