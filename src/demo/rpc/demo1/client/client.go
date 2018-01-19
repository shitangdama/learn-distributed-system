package main

import (
	"fmt"
	"net/rpc"
	"os"
)

type Args struct {
	I, J int
}

type DivResult struct {
	Quo, Rem int
}

func main() {
	// 客户端先用rpc.DialHTTP和RPC服务器进行一个链接(协议必须匹配).
	pClient, err := rpc.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err : %s\n", err.Error())
		return
	}

	// 同步RPC
	var multResult int
	err = pClient.Call("MyMethod.Mult", &Args{2, 7}, &multResult)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err : %s\n", err.Error())
		return
	}
	fmt.Println("Mult返回", multResult)

	// 异步RPC
	var divResult DivResult
	pCall := pClient.Go("MyMethod.Div", &Args{5, 2}, &divResult, nil)
	if pCall != nil {
		if replyCall, ok := <-pCall.Done; ok {
			fmt.Println("Div异步收到返回", replyCall)
			fmt.Println("Div异步收到返回", divResult.Quo, divResult.Rem)
		}
	}
}