package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// сюда писать код
func SingleHash(in, out chan interface{}) {
	//wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}

	// crc32(data) + "~" + md5(data)
	// распараллеливаем вызов функций
	util(in, out, func(data string) string {
		// записываем полученные хэш-значения функции в 2 канала и возвращаем их
		// функция долго работает
		//wg.Add(2)
		chan1 := make(chan string)
		chan2 := make(chan string)
		// закидываем данные в первый канал и считаем crc32
		go func(data string) {
			defer close(chan1)
			chan1 <- DataSignerCrc32(data)
		}(data)

		// закидываем данные во второй канал и считаем md5(одновременно 1 такая ф-ия)
		go func(data string) {
			defer close(chan2)
			mutex.Lock()
			data_ := DataSignerMd5(data)
			mutex.Unlock()
			chan2 <- data_
		}(data)
		return fmt.Sprintf("%s~%s", <-chan1, <-chan2)
	})
}

func MultiHash(in, out chan interface{}) {
	util(in, out, func(data string) string {
		array_ := make([]string, 6)
		wg := &sync.WaitGroup{}
		mutex := &sync.Mutex{}
		for i := 0; i < 6; i++ {
			wg.Add(1)
			go func(thread int) {
				defer wg.Done()
				hash := DataSignerCrc32(fmt.Sprintf("%d%s", thread, data))
				mutex.Lock()
				array_[thread] = hash
				mutex.Unlock()
			}(i)
		}
		wg.Wait()
		return strings.Join(array_, "")
	})
}

func CombineResults(in, out chan interface{}) {
	var ans []string
	for i := range in {
		ans = append(ans, fmt.Sprintf("%v", i))
	}
	sort.Strings(ans)
	out <- strings.Join(ans, "_")
}

func util(in, out chan interface{}, str_func func(string) string) {
	wg := &sync.WaitGroup{}
	// проходимся по всем данным из входного потока в waitGroup и записываем данные в выходной поток
	for i := range in {
		wg.Add(1)
		go func(str string) {
			defer wg.Done()
			out <- str_func(str)
		}(fmt.Sprintf("%v", i))
	}
	wg.Wait()
}

// функция проходит по массиву job-ов, добавляет
func ExecutePipeline(jobs ...job) {
	var in, out chan interface{}
	in = make(chan interface{})

	wg := &sync.WaitGroup{}
	for _, job_ := range jobs {
		out = make(chan interface{}, MaxInputDataLen)
		wg.Add(1)
		go func(j job, in_, out_ chan interface{}) {
			defer close(out_)
			defer wg.Done()
			j(in_, out_)
		}(job_, in, out)
		in = out
	}
	wg.Wait()
}
