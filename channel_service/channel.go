package channel_service

import (
	"sync"
	"time"
)

//没有缓冲的管道，对于接收方会一直阻塞直至管道中被加入数据
func SyncWithNoBufferChannel(){
	c := make(chan string)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		c <- `hello Golang`
	}()
	go func() {
		defer wg.Done()
		println(`Message: `+ <-c)
	}()
	wg.Wait()
}
func SyncWithHaveBufChannel(){
		c := make(chan string, 2)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			c <- `world`
			c <- `lyh`
		}()
		go func() {
			defer wg.Done()
			time.Sleep(time.Second * 1)
			println(`hello: `+ <-c)
			println(`作者: `+ <-c)
		}()
		wg.Wait()
}