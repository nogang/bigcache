package main

import (
	"fmt"
	"time"
	"github.com/allegro/bigcache"
	"runtime"
	"sync"
	"github.com/hashicorp/golang-lru"
	"math/rand"
	"math"
)
var maxGoroutine int
func main(){
	maxGoroutine = 1
	fmt.Println("test world")
	//memBigCacheTest()
	memALUTest()
	printMem()
}
var wait sync.WaitGroup
var maxAlloc float64

func memALUTest(){
	printMem()
	cache, _ := lru.New(80960)
	var wait sync.WaitGroup
	wait.Add(1)
	for goCount := 0 ; goCount < 100 ; goCount++ {
		go func() {
			rand.Seed(time.Now().UTC().UnixNano())
			for {
				for i := 0; i < 20000; i++ {
					value := make([]byte, 1024*8)
					cache.Add(fmt.Sprintf("%d", rand.Int()), value)
				}
				printMem()
			}
		} ()
	}

	wait.Wait()
		/*
	for i:= 0; i < 1000000 ; i++ {
		_,ok := cache.Get(fmt.Sprintf("%d",i))
		if !ok {
			fmt.Println("not")
		}

	}
*/


}

func memBigCacheTest(){

	cache, _ := bigcache.NewBigCache(bigcache.Config{
		Shards:             256,
		LifeWindow:         10 * time.Minute,
		MaxEntriesInWindow: 10000,
		MaxEntrySize:       256,
		Verbose:            false,
	})
	printMem()

	wait.Add(maxGoroutine)
	start := time.Now()
	for t := 0 ; t < maxGoroutine ; t++ {
		go testBigCacheFunc(cache, t)
	}
	fmt.Printf("gorou : %d \n", runtime.NumGoroutine())
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	fmt.Printf("%s", buf)

	wait.Wait()
	elapsedTime := time.Since(start)
	fmt.Printf("%s",elapsedTime/10000000)
	printMem()
}

func testBigCacheFunc(bc *bigcache.BigCache,start int){
	defer wait.Done()
	value := make([]byte, 256)
	for i:= 0; i < 1000000/maxGoroutine ; i++ {
		e := bc.Set(fmt.Sprintf("%d",start + i),value)
		if e != nil {
			fmt.Errorf("err %s", e)
		}
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
var count int
func printMem(){
	count++
	fmt.Printf("")
	var  m runtime.MemStats
	runtime.ReadMemStats(&m)
	maxAlloc = math.Max(maxAlloc, float64(m.Alloc))
	fmt.Printf("maxAllo = %v MiB", bToMb(uint64(maxAlloc)))
	fmt.Printf("\tAlloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v", m.NumGC)
	fmt.Printf("\tfree = %v MiB\n", bToMb(m.Frees))

}