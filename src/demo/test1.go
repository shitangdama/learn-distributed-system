// 这一个测试，编写同步锁的相关
// go mutex是互斥锁，只有Lock和Unlock两个方法，在这两个方法之间的代码不能被多个goroutins同时调用到。
package demo

import (
	"fmt"
	"sync"
	"time"
)

var m *sync.Mutex

func main() {

	m = new(sync.Mutex)
	go lockPrint(1)
	lockPrint(2)
	time.Sleep(time.Second)
	fmt.Printf("%s\n", "exit!")
}

func lockPrint(i int) {
	println(i, "lock start")
	m.Lock()
	println(i, "in lock")
	time.Sleep(3 * time.Second)
	m.Unlock()
	println(i, "unlock")
}

// 2 lock start

// 2 in lock

// 1 lock start

// 2 unlock

// 1 in lock

// exit!