package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	ch2 := make(chan string)
	go func() {
		//for {
		time.Sleep(time.Millisecond * 500)
		ch <- "Прошло пол-секунды"
		close(ch)
		//}
	}()

	go func() {
		//for {
		time.Sleep(time.Millisecond * 2000)
		ch2 <- "Прошло 2 секунды"
		close(ch2)
		//}
	}()
	// nio из нескольких каналов
	select {
	case msg := <-ch:
		fmt.Println(msg)
	case msg := <-ch2:
		fmt.Println(msg)
	}
}
