package main

import (
	"fmt"
	"os"
	"regexp"
	"sync"
)

//RingBuffer - Кольцевой буфер
type RingBuffer struct {
	full   int
	empty  int
	busy   int
	buffer []string
	head   int
	tail   int
}

func newRingBuffer(buffSize int) *RingBuffer {
	rb := RingBuffer{
		full:   0,
		empty:  buffSize,
		busy:   1,
		buffer: make([]string, buffSize),
		head:   0,
		tail:   0}
	return &rb
}

//V(S)
func up(semaphore *int) {
	*semaphore = *semaphore + 1
}

//P(S)
func down(semaphore *int) {
	if *semaphore > 0 {
		*semaphore = *semaphore - 1
	}
}

func producer(buffer *RingBuffer, wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		fmt.Print("Введите любой символ (введите kill для выхода) : ")
		var item string
		fmt.Scan(&item)

		down(&buffer.empty)
		down(&buffer.busy)

		if item != "kill" {
			buffer.buffer[buffer.head] = item
			buffer.head = (buffer.head + 1) % len(buffer.buffer)
		}

		up(&buffer.busy)
		up(&buffer.full)

		if item == "kill" {
			wg.Add(-3)
			break
		}
	}
}

//ловит буквы английского алфавита из буфера
func consumer1(buffer *RingBuffer, wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		down(&buffer.full)
		down(&buffer.busy)

		var item *string

		if (buffer.head-buffer.tail) == 1 || (buffer.tail-buffer.head) == (len(buffer.buffer)-1) {

			item = &buffer.buffer[buffer.tail]

			match, _ := regexp.MatchString("^[a-zA-Z]$", *item)

			if match {
				buffer.tail = (buffer.tail + 1) % len(buffer.buffer)
			} else {
				item = nil
			}
		}

		up(&buffer.busy)
		up(&buffer.empty)

		if item != nil {
			saveTextToFile("./consumer1.txt", *item+"\n")
		}
	}
}

//ловит цифры из буфера
func consumer2(buffer *RingBuffer, wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		down(&buffer.full)
		down(&buffer.busy)

		var item *string

		if (buffer.head-buffer.tail) == 1 || (buffer.tail-buffer.head) == (len(buffer.buffer)-1) {

			item = &buffer.buffer[buffer.tail]

			match, _ := regexp.MatchString("^[0-9]$", *item)

			if match {
				buffer.tail = (buffer.tail + 1) % len(buffer.buffer)
			} else {
				item = nil
			}
		}

		up(&buffer.busy)
		up(&buffer.empty)

		if item != nil {
			saveTextToFile("./consumer2.txt", *item+"\n")
		}
	}
}

//ловит остальные символы из буфера
func consumer3(buffer *RingBuffer, wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		down(&buffer.full)
		down(&buffer.busy)

		var item *string

		if (buffer.head-buffer.tail) == 1 || (buffer.tail-buffer.head) == (len(buffer.buffer)-1) {

			item = &buffer.buffer[buffer.tail]

			match, _ := regexp.MatchString("^[\\W]$", *item)

			if match {
				buffer.tail = (buffer.tail + 1) % len(buffer.buffer)
			} else {
				item = nil
			}
		}

		up(&buffer.busy)
		up(&buffer.empty)

		if item != nil {
			saveTextToFile("./consumer3.txt", *item+"\n")
		}
	}
}

func main() {

	buffer := newRingBuffer(5)

	var wg sync.WaitGroup

	wg.Add(1)
	go producer(buffer, &wg)

	wg.Add(1)
	go consumer1(buffer, &wg)

	wg.Add(1)
	go consumer2(buffer, &wg)

	wg.Add(1)
	go consumer3(buffer, &wg)

	wg.Wait()

	fmt.Println()
	fmt.Println(buffer)

}

func saveTextToFile(filename, text string) {

	file, err := os.OpenFile(filename, os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		panic(err)
	}

}
