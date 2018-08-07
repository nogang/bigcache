package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/allegro/bigcache"


			"github.com/hashicorp/golang-lru"
	"runtime"
	"math/rand"
)

const maxEntrySize = 256
const mapSize = 100
const maxGoroutine = 5
const multiCache = 2



////////////////cache size별 get time////////////

////////////////cache size별 read time////////////

type Cache interface{
	Add(key, value interface{}) (evicted bool)
	Get(key interface{}) (value interface{}, ok bool)
}

type lruCache struct {
	lru *lru.Cache
}

func (cache *lruCache) Add(key, value interface{}) (evicted bool) {
	return cache.lru.Add(key, value)
}

func (cache *lruCache) Get(key interface{}) (value interface{}, ok bool) {
	return cache.lru.Get(key)
}

type bigCache struct {
	big *bigcache.BigCache
}

func (cache *bigCache) Add(key, value interface{}) (evicted bool) {
	k, kok := key.(string)
	v, vok := value.([]byte)
	if kok && vok {
		cache.big.Set(k, v)
	}

	return false
}

func (cache *bigCache) Get(key interface{}) (value interface{}, ok bool) {
	k, ok := key.(string)
	if ok {
		ret, error := cache.big.Get(k)
		return ret, error == nil

	}
	return nil, false
}

type BM struct{
	name		string
	cacheSize	int
	inDataSize 	int
}



type TestFunc func(b *testing.B, cache Cache, bm BM)

func BenchmarkCacheSingleGetTest(b *testing.B){
	benchCacheTest(b, singleGetTestFunc)
}
func BenchmarkCacheSingleAddTest(b *testing.B){
	benchCacheTest(b, singleAddTestFunc)
}

func BenchmarkCacheParellalGetTest(b *testing.B){
	benchCacheTest(b, parallelGetTestFunc)
}

func BenchmarkCacheParellalAddTest(b *testing.B){
	benchCacheTest(b, parallelAddTestFunc)
}

func benchCacheTest(b *testing.B, tf TestFunc){
	benchmarks := []BM{}
	cacheName := []string{"lru", "bigCache"}

	for i := 0 ; i < len(cacheName) ; i++ {
		for cacheSize := 1000 ; cacheSize <= 10000000 ; cacheSize *= 10 {
		//for cacheSize := 1000 ; cacheSize <= 1000 ; cacheSize *= 10 {
			benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize/10})
			benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize/100})
			benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize/1000})
		}
	}

	var cache Cache
	for _, bm := range benchmarks {
		cache, _ = newCache(bm.name, bm.cacheSize)

		for i := 0; i < bm.inDataSize; i++ {
			cache.Add(key(i), value())
		}

		tf(b, cache, bm)
	}
}

var singleAddTestFunc = func(b *testing.B, cache Cache, bm BM){
	testName := fmt.Sprintf("%s,cacheSize(count):%d,inData(count):%d",bm.name, bm.cacheSize, bm.inDataSize)
	b.Run(testName, func(b *testing.B) {
		for i:= 0 ; i < b.N ; i++{
			cache.Add(key(i+bm.inDataSize), value())
		}
	})
}


var singleGetTestFunc = func(b *testing.B, cache Cache, bm BM){
	testName := fmt.Sprintf("%s,cacheSize(count):%d,inData(count):%d",bm.name, bm.cacheSize, bm.inDataSize)
	b.Run(testName, func(b *testing.B) {
		for i:= 0 ; i < b.N ; i++{
			cache.Get(key(i%bm.inDataSize))
		}
	})
}

var parallelAddTestFunc = func(b *testing.B, cache Cache, bm BM){
	testName := fmt.Sprintf("%s,cacheSize(count):%d,inData(count):%d,goRoutine:%d",bm.name, bm.cacheSize, bm.inDataSize, maxGoroutine)
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.Run(testName, func(b *testing.B){
		b.RunParallel(func(pb *testing.PB) {
			id := rand.Intn(maxGoroutine * 1000)
			counter := 0
			for pb.Next() {
				cache.Add(parallelKey(id, counter), value())
				counter++
			}
		})
	})
}

var parallelGetTestFunc = func(b *testing.B, cache Cache, bm BM){
	testName := fmt.Sprintf("%s,cacheSize(count):%d,inData(count):%d,goRoutine:%d",bm.name, bm.cacheSize, bm.inDataSize, maxGoroutine)
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.Run(testName, func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			counter := 0
			for pb.Next() {
				cache.Get(key(counter % bm.inDataSize))
				counter++
			}
		})
	})
}

func newCache(cacheName string, size int ) (Cache, error){
	switch cacheName {
	case "lru":
		return lru.New(size)
	case "bigCache" :
		return initBigCache(size, 1024), nil
	}
	return nil, nil
}

/////랜덤 테스트 구현/////
/*


func BenchmarkCacheAddTest(b *testing.B){
	benchmarks := []BM{}
	cacheName := []string{"lru", "bigCache"}

	for i := 0 ; i < len(cacheName) ; i++ {
		for cacheSize := 100 ; cacheSize <= 100000000 ; cacheSize *= 10 {
			benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize/10})
			benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize/100})
			benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize/1000})
		}
	}

	var cache Cache
	for _, bm := range benchmarks {
		cache, _ = newCache(bm.name, bm.cacheSize)

		for i := 0; i < bm.inDataSize; i++ {
			cache.Add(key(i), value())
		}

		testName := fmt.Sprintf("%s,cacheSize(count):%d,inData(count):%d",bm.name, bm.cacheSize, bm.inDataSize)
		b.Run(testName, func(b *testing.B) {
			for i:= 0 ; i < b.N ; i++{
				cache.Add(key(i+bm.inDataSize), value())
			}
		})
	}
}

func BenchmarkGoCacheSetParallel(b *testing.B) {
	c := cache.New(5*time.Minute, 10*time.Minute)

	rand.Seed(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(maxGoroutine * 1000)
		counter := 0
		for pb.Next() {
			c.Set(parallelKey(id, counter), value(),cache.DefaultExpiration)
			counter = counter + 1
		}
	})
}
*/
/*
func BenchmarkPeekTest(b *testing.B){
	cache, _ := lru.New(b.N * multiCache)
	for i := 0; i < b.N; i++ {
		cache.Add(key(i), value())
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Peek(key(counter))
			counter = counter + 1
		}
	})
}

func BenchmarkLRUCacheGetParallel2(b *testing.B) {
	cache, _ := lru.New(b.N * multiCache)
	for i := 0; i < b.N; i++ {
		cache.Add(key(i), value())
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Get(key(counter))
			counter = counter + 1
		}
	})
}



func BenchmarkSyncMapSet(b *testing.B){
	var m sync.Map
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			m.Store(key(counter), value())
			counter++
		}
	})
}
func BenchmarkSyncMapGet(b *testing.B){
	var m sync.Map
	for i := 0; i < b.N; i++ {
		m.Store(key(i), value())
	}

	hitCount := 0
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			_, ok := m.Load(key(counter))
			if ok {
				hitCount++
			}
			counter++
		}
	})
}

func benchmarkBigCacheGetShard(b *testing.B, shard int) {
	//b.StopTimer()
	cache := initBigCache(b.N * multiCache, shard)
	rand.Seed(time.Now().Unix())
	for i := 0; i < b.N; i++ {
		cache.Set(key(i), value())
	}
	//b.StartTimer()
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Get(key(counter))
			counter = counter + 1
		}
	})
}


func BenchmarkGoCacheSetParallel(b *testing.B) {
	c := cache.New(5*time.Minute, 10*time.Minute)

	rand.Seed(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(maxGoroutine * 1000)
		counter := 0
		for pb.Next() {
			c.Set(parallelKey(id, counter), value(),cache.DefaultExpiration)
			counter = counter + 1
		}
	})
}

func BenchmarkGoCacheGetParallel(b *testing.B) {
	c := cache.New(5*time.Minute, 10*time.Minute)
	for i := 0; i < b.N; i++ {
		c.Set(key(i), value(), cache.DefaultExpiration)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			c.Get(key(counter))
			counter = counter + 1
		}
	})
}


func BenchmarkLRUCacheSetParallel(b *testing.B) {
	cache, _ := lru.New(b.N * multiCache)
	rand.Seed(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(maxGoroutine * 1000)
		counter := 0
		for pb.Next() {
			cache.Add(parallelKey(id, counter), value())
			counter = counter + 1
		}
	})
}

func BenchmarkLRUCacheGetParallel(b *testing.B) {
	cache, _ := lru.New(b.N * multiCache)
	for i := 0; i < b.N; i++ {
		cache.Add(key(i), value())
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Get(key(counter))
			counter = counter + 1
		}
	})
}

func BenchmarkFreeCacheSetParallel(b *testing.B) {
	cache := freecache.NewCache(b.N * maxEntrySize * multiCache)
	rand.Seed(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(maxGoroutine * 1000)
		counter := 0
		for pb.Next() {
			cache.Set([]byte(parallelKey(id, counter)), value(), 0)
			counter = counter + 1
		}
	})
}

func BenchmarkFreeCacheGetParallel(b *testing.B) {
	cache := freecache.NewCache(b.N * maxEntrySize * multiCache)
	for i := 0; i < b.N; i++ {
		cache.Set([]byte(key(i)), value(), 0)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Get([]byte(key(counter)))
			counter = counter + 1
		}
	})
}
*/
/*
func BenchmarkBigCacheSetParallel1024(b *testing.B) {
	benchmarkBigCacheSetParallel(b, 1024)
}
func BenchmarkBigCacheGetParallel1024Test(b *testing.B) {
	benchmarkBigCacheGetShard(b,1024)
}
*/
/*
func BenchmarkBigCacheGetShard64Test(b *testing.B) {
	benchmarkBigCacheGetShard(b,64)
}
*/

/*
func BenchmarkBigCacheGetShard2048Test(b *testing.B) {
	benchmarkBigCacheGetShard(b,2048)
}
*/

/*
func benchmarkBigCacheSetParallel(b *testing.B, shard int) {
	cache := initBigCache(b.N * multiCache, shard)
	rand.Seed(time.Now().Unix())
	b.SetParallelism(maxGoroutine)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(maxGoroutine * 1000)
		counter := 0
		for pb.Next() {
			cache.Set(parallelKey(id, counter), value())
			counter = counter + 1
		}
	})
}
*/
/*
func BenchmarkBigCacheSetParallel64(b *testing.B) {
	benchmarkBigCacheSetParallel(b, 64)
}
*/

/*

func BenchmarkBigCacheSetParallel2048(b *testing.B) {
	benchmarkBigCacheSetParallel(b, 2048)
}

*/
func initBigCache(entriesInWindow int, shards int) *bigCache {
	cache, _ := bigcache.NewBigCache(bigcache.Config{
		Shards:             shards,
		LifeWindow:         10 * time.Minute,
		MaxEntriesInWindow: entriesInWindow,
		MaxEntrySize:       maxEntrySize,
		Verbose:            false,
	})

	return &bigCache{cache}
}


/*

func BenchmarkConcurrentMapGetParallel(b *testing.B) {
	b.StopTimer()
	var m sync.Map
	for i := 0; i < b.N; i++ {
		m.Store(key(i), value())
	}

	b.StartTimer()
	hitCount := 0

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		for pb.Next() {
			_, ok := m.Load(key(id))
			if ok {
				hitCount++
			}
		}
	})
}
*/

/*
func BenchmarkBigCacheTest(b *testing.B) {
	for i:= 0 ; i < 1 ; i++ {
		ecount := 0
		cache := initBigCache(1000000000)
		for i := 0; i < 100000000; i++ {
			cache.Set(key(i), value())
		}
		fmt.Println("ore : %d", ecount)
		for i := 0; i < 100000000; i++ {
			_, e := cache.Get(key(i))
			if e != nil{
				ecount++
			} else {
				//fmt.Println(v)
			}
		}
		fmt.Println("result : %d", ecount)
	}


}

func BenchmarkBigCacheSet10000Level2(b *testing.B) {
	for k :=0; k < b.N; k++ {
		cache := initBigCache(10000 * 10)
		for i := 0; i < 10000; i++ {
			cache.Set(key(i), value())
		}
	}
}

func BenchmarkBigCacheSet10000Level3(b *testing.B) {
	for k :=0; k < b.N; k++ {
		cache := initBigCache(10000 * 50)
		for i := 0; i < 10000; i++ {
			cache.Set(key(i), value())
		}
	}
}

func BenchmarkBigCacheSet10000Level4(b *testing.B) {
	for k :=0; k < b.N; k++ {
		cache := initBigCache(10000 * 100)
		for i := 0; i < 10000; i++ {
			cache.Set(key(i), value())
		}
	}
}

func BenchmarkLRUCacheSet(b *testing.B) {
	cache, _ := lru.New(b.N)
	for i := 0; i < b.N; i++ {
		cache.Add(key(i), value())
	}
}

func benchmarkMapSetSize(b *testing.B) {
	b.StopTimer()
	m := make(map[string][]byte)
	b.StartTimer()
	for i := 0; i < mapSize; i++ {
		m[key(i)] = value()
	}
}

func BenchmarkMapSet(b *testing.B){
	for i := 0; i < b.N; i++ {
		benchmarkMapSetSize(b)
	}
}

func benchmarkConcurrentMapSet() {
	var m sync.Map
	for i := 0; i < mapSize; i++ {
		m.Store(key(i), value())
	}
}

func BenchmarkConcurrentMapSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkConcurrentMapSet()
	}
}

func benchmarkFreeCacheSet(b *testing.B) {
	b.StopTimer()
	cache := freecache.NewCache(mapSize * maxEntrySize)
	b.StartTimer()
	for i := 0; i < mapSize; i++ {
		cache.Set([]byte(key(i)), value(), 0)
	}
}

func BenchmarkFreeCacheSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkFreeCacheSet(b)
	}
}

func benchmarkBigCacheSet(b *testing.B) {
	b.StopTimer()
	cache := initBigCache(mapSize)
	b.StartTimer()
	for i := 0; i < mapSize; i++ {
		cache.Set(key(i), value())
	}
}

func BenchmarkBigCacheSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkBigCacheSet(b)
	}
}

////Get////
func BenchmarkMapGet(b *testing.B) {
	b.StopTimer()
	m := make(map[string][]byte)
	for i := 0; i < mapSize; i++ {
		m[key(i)] = value()
	}

	b.StartTimer()
	hitCount := 0
	for re:= 0 ; re < b.N ; re++ {
		for i := 0; i < mapSize; i++ {
			if m[key(i)] != nil {
				hitCount++
			}
		}
	}
}

func BenchmarkConcurrentMapGet(b *testing.B) {
	b.StopTimer()
	var m sync.Map
	for i := 0; i < mapSize; i++ {
		m.Store(key(i), value())
	}

	b.StartTimer()
	hitCounter := 0
	for re := 0; re < b.N; re++ {
		for i := 0; i < mapSize; i++ {
			_, ok := m.Load(key(i))
			if ok {
				hitCounter++
			}
		}
	}
}

func BenchmarkFreeCacheGet(b *testing.B) {
	b.StopTimer()
	cache := freecache.NewCache(mapSize * maxEntrySize)
	for i := 0; i < mapSize; i++ {
		cache.Set([]byte(key(i)), value(), 0)
	}

	b.StartTimer()
	for re := 0; re < b.N; re++ {
		for i := 0; i < mapSize; i++ {
			cache.Get([]byte(key(i)))
		}
	}
}

func BenchmarkBigCacheGet(b *testing.B) {
	b.StopTimer()
	cache := initBigCache(mapSize)
	for i := 0; i < mapSize; i++ {
		cache.Set(key(i), value())
	}

	b.StartTimer()
	for re := 0; re < b.N; re++ {
		for i := 0; i < mapSize; i++ {
			cache.Get(key(i))
		}
	}
}

////Set Parallel////

func BenchmarkBigCacheSetParallel(b *testing.B) {
	cache := initBigCache(b.N)
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set(parallelKey(id, counter), value())
			counter = counter + 1
		}
	})
}

func BenchmarkFreeCacheSetParallel(b *testing.B) {
	cache := freecache.NewCache(b.N * maxEntrySize)
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set([]byte(parallelKey(id, counter)), value(), 0)
			counter = counter + 1
		}
	})
}

func BenchmarkConcurrentMapSetParallel(b *testing.B) {
	var m sync.Map

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		for pb.Next() {
			m.Store(key(id), value())
		}
	})
}

func BenchmarkBigCacheGetParallel(b *testing.B) {
	b.StopTimer()
	cache := initBigCache(b.N)
	for i := 0; i < b.N; i++ {
		cache.Set(key(i), value())
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Get(key(counter))
			counter = counter + 1
		}
	})
}


////Get parallel////
func BenchmarkFreeCacheGetParallel(b *testing.B) {
	b.StopTimer()
	cache := freecache.NewCache(b.N * maxEntrySize)
	for i := 0; i < b.N; i++ {
		cache.Set([]byte(key(i)), value(), 0)
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Get([]byte(key(counter)))
			counter = counter + 1
		}
	})
}

func BenchmarkConcurrentMapGetParallel(b *testing.B) {
	b.StopTimer()
	var m sync.Map
	for i := 0; i < b.N; i++ {
		m.Store(key(i), value())
	}

	b.StartTimer()
	hitCount := 0

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		for pb.Next() {
			_, ok := m.Load(key(id))
			if ok {
				hitCount++
			}
		}
	})
}

*/
func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}
func value() []byte {
	return make([]byte, 100)
}


func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}
