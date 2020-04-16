package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("run...")
	// sum := addAfter(8, 100, 3*time.Second)
	// fmt.Println(<-sum)

	// sum := addByTick(1, 100, 5*time.Second)
	// for v := range sum {
	// 	fmt.Println(v)
	// }

	ctx, cancel := context.WithCancel(context.Background())
	arr := makeArray(ctx, 1, 9999999)
	// var arr = make([]int, 9999999)
	PrintMemUsage()
	// for v := range arr {
	// 	// fmt.Println(v)
	// 	v++
	// }

	for i := 0; i < 5; i++ {
		fmt.Println(<-arr)
	}
	for i := 0; i < 5; i++ {
		fmt.Println(<-arr)
	}
	cancel()
	for i := 0; i < 5; i++ {
		fmt.Println(<-arr)
	}

	PrintMemUsage()

}

func addAfter(a, b int, seconds time.Duration) chan int {
	ch := make(chan int, 0)
	time.AfterFunc(seconds, func() {
		ch <- a + b
	})
	return ch
}
func addByTick(base, offset int, seconds time.Duration) chan int {
	ch := make(chan int, 0)
	ticker := time.NewTicker(1 * time.Second)
	quit := time.After(seconds)
	tBase := base
	go func() {
		for {
			select {
			case <-ticker.C:
				tBase += offset
				ch <- tBase
			case <-quit:
				ticker.Stop()
				close(ch)
			}
		}
	}()
	return ch
}

func makeArray(ctx context.Context, from, to int) chan int {
	ch := make(chan int, 0)
	go func() {
	tag:
		for i := from; i < to; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("ctx.Done")
				// close(ch)
				break tag
			case ch <- i:

			}
		}
		close(ch)
	}()
	return ch
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
