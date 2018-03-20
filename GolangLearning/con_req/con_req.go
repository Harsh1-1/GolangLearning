package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/trace"
	"time"
)

func send_request(s string) {
	resp, _ := http.Get(s)
	fmt.Println(s, "\n", resp, "\n")
}
func main() {

	runtime.GOMAXPROCS(4)
	// for tracing
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	var s []string
	fileHandle, _ := os.Open("random_domains.txt")
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)
	for fileScanner.Scan() {
		// fmt.Println(fileScanner.Text())
		s = append(s, "http://"+fileScanner.Text())
	}
	// fmt.Printf("%v\n", s)
	for _, url := range s {
		go send_request(url)
	}

	time.Sleep(30 * time.Second)

	// stopping trace
	trace.Stop()
	if err != nil {
		panic(err)
	}
}
