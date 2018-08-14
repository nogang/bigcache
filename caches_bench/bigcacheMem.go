package main

import (
	"fmt"
	"time"
	"github.com/allegro/bigcache"
	"runtime"
	"sync"
)
var maxGoroutine int
func main(){
	maxGoroutine = 5000
	fmt.Println("test world")
	memTest()
}
var wait sync.WaitGroup
func memTest(){

	cache, _ := bigcache.NewBigCache(bigcache.Config{
		Shards:             256,
		LifeWindow:         10 * time.Minute,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       256,
		Verbose:            false,
	})
	printMem()

	wait.Add(maxGoroutine)
	for t := 0 ; t < maxGoroutine ; t++ {
		go testFunc(cache, t)
	}

	wait.Wait()
	printMem()
}

func testFunc(bc *bigcache.BigCache,start int){
	defer wait.Done()
	value := make([]byte, 256)
	for i:= 0; i < 10000000/maxGoroutine ; i++ {
		e := bc.Set(fmt.Sprintf("%d",start + i),value)
		if e != nil {
			fmt.Errorf("err %s", e)
		}
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func printMem(){
	var  m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}