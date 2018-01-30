package main

import (
	//	"fmt"
	"time"
)

func main() {
	//	ticker := time.NewTimer(time.Second * 5)
	//	go func() {
	//		for _ = range ticker.C { //
	//			fmt.Println("aa")
	//		}
	//	}()
	msg1 := make(chan interface{})
	demo(msg1)
}

//实现了定时任务
func demo(input chan interface{}) {
	t1 := time.NewTimer(time.Second * 5)
	t2 := time.NewTimer(time.Second * 10)
	msg := make(chan interface{})

	for {
		select {
		case msg <- input:
			println(msg)

		case <-t1.C:
			println("5s timer")
			t1.Reset(time.Second * 5)

		case <-t2.C:
			println("10s timer")
			t2.Reset(time.Second * 10) //重置一个10s的定时器
		}
	}
}
// go demo(input)

// func demo(input chan interface{}) {
//     for {
//         select {
//         case msg <- input:
//             println(msg)

//         case <-time.After(time.Second * 5):
//             println("5s timer")

//         case <-time.After(time.Second * 10):
//             println("10s timer")
//         }
//     }
// }
// 写出上面这段程序的目的是从 input channel 持续接收消息加以处理，同时希望每过5秒钟和每过10秒钟就分别执行一个定时任务。
// 但是当你执行这段程序的时候，只要 input channel 中的消息来得足够快，永不间断，你会发现启动的两个定时任务都永远不会执行；
// 即使没有消息到来，第二个10s的定时器也是永远不会执行的。
// 原因就是 select 每次执行都会重新执行 case 条件语句，并重新注册到 select 中，因此这两个定时任务在每次执行 select 的时候，
// 都是启动了一个新的从头开始计时的 Timer 对象，所以这两个定时任务永远不会执行。