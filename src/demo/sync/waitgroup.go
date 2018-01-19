package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var urls = []string{
		"http://www.sohu.com/",
		"http://www.baidu.com/",
		"http://www.feixun.tv/",
	}
	for _, url := range urls {
		// Increment the WaitGroup counter.
		// 递增 WaitGroup 计数器。
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		// 启动一个Go程来取回URL。
		go func(url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Fetch the URL.
			// 取回URL
			res, _ := http.Get(url)
			fmt.Println(res)
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	// 等待所有的HTTP取回操作完成。
	wg.Wait()
}