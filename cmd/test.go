package main

import (
	"fmt"

	"github.com/gammazero/workerpool"
)

func test() {
	wp := workerpool.New(5)
	for i := 0; i < 100; i++ {
		j := i
		wp.Submit(func() {
			fmt.Println(j)
		})
	}

	wp.StopWait()
	fmt.Println("Im done")
}
