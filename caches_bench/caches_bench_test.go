package main

import (
	"fmt"
	"testing"
	"time"
	"github.com/allegro/bigcache"
	"github.com/hashicorp/golang-lru"
	"runtime"
	"math/rand"
	"github.com/coocood/freecache"
	GOCACHE "github.com/patrickmn/go-cache"
	"sync"
)

const maxEntrySize = 256
const mapSize = 100
const maxGoroutine = 5
const multiCache = 2


const LRU = "LRU"
const BigCache = "BigCache"
const FreeCache = "FreeCache"
const GoCache = "GoCache"
const SyncMap = "SyncMap"

////////////////cache size별 get time////////////

////////////////cache size별 read time////////////

type Cache interface{
	Add(key, value interface{}) (ok bool)
	Get(key interface{}) (value interface{}, ok bool)
}

type lruCache struct {
	lru *lru.Cache
}

func (cache *lruCache) Add(key, value interface{}) (ok bool) {
	cache.lru.Add(key, value)
	return true
}

func (cache *lruCache) Get(key interface{}) (value interface{}, ok bool) {
	return cache.lru.Get(key)
}

type bigCache struct {
	big *bigcache.BigCache
}

func (cache *bigCache) Add(key, value interface{}) (ok bool) {
	k, kok := key.(string)
	v, vok := value.([]byte)
	if ok = kok && vok ; ok {
		cache.big.Set(k, v)
	}
	return
}

func (cache *bigCache) Get(key interface{}) (value interface{}, ok bool) {
	k, ok := key.(string)
	if ok {
		ret, error := cache.big.Get(k)
		return ret, error == nil
	}
	return nil, false
}

type freeCache struct{
	free *freecache.Cache
}

func (cache *freeCache)Add(key, value interface{}) (ok bool) {
	k, kok := key.([]byte)
	v, vok := value.([]byte)

	if ok = kok && vok ; ok {
		cache.free.Set(k, v,0)
	}
	return
}

func (cache *freeCache) Get(key interface{}) (value interface{}, ok bool) {
	k, ok := key.([]byte)
	if ok {
		ret, error := cache.free.Get(k)
		return ret, error == nil
	}
	return nil, false
}

type goCache struct {
	cache *GOCACHE.Cache
}

func (gc *goCache)Add(key, value interface{}) (ok bool) {
	k, ok := key.(string)

	if ok {
		gc.cache.Set(k,value,1000)
	}
	return
}

func (gc *goCache) Get(key interface{}) (value interface{}, ok bool) {
	k, ok := key.(string)
	if ok {
		return gc.cache.Get(k)
	}
	return nil, false
}

type syncMap struct {
	cache sync.Map
}

func (sm *syncMap)Add(key, value interface{}) (ok bool) {
	sm.cache.Store(key, value)
	return true
}

func (sm *syncMap)Get(key interface{}) (value interface{}, ok bool) {
	return sm.cache.Load(key)
}

func newCache(cacheName string, size int ) (Cache, error){
	switch cacheName {
	case LRU:
		return lru.New(size)
	case BigCache:
		return initBigCache(size, 1024), nil
	case FreeCache:
		return initFreeCache(size), nil
	case GoCache:
		return initGoCache(), nil
	case SyncMap:
		return initSyncMap(), nil

	}
	return nil, nil
}

type BM struct{
	name		string
	cacheSize	int
	inDataSize 	int
}

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

func initFreeCache(size int) *freeCache{
	return &freeCache{free:freecache.NewCache(size * maxEntrySize)}
}

func initGoCache() *goCache{
	return &goCache{cache:GOCACHE.New(5*time.Minute, 10*time.Minute)}
}

func initSyncMap() *syncMap{
	var m sync.Map
	return &syncMap{cache:m}
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
	cacheName := []string{LRU, BigCache, FreeCache, GoCache}

	for i := 0 ; i < len(cacheName) ; i++ {
		for cacheSize := 1000 ; cacheSize <= 10000000 ; cacheSize *= 10 {
		//for cacheSize := 1000 ; cacheSize <= 1000 ; cacheSize *= 10 {
			if cacheName[i] == GoCache || cacheName[i] == SyncMap {
				benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize})
			} else {
				benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize/10})
				benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize/100})
				benchmarks = append(benchmarks,BM{name:cacheName[i],cacheSize:cacheSize,inDataSize:cacheSize/1000})
			}
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

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}
func value() []byte {
	return make([]byte, 100)
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}


/////랜덤 테스트 구현/////
