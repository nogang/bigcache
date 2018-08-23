package main

import (
	"testing"
	"github.com/hashicorp/golang-lru"
	"fmt"
)

type BM_DATASIZE struct{
	datasize int
	data 	interface{}
}

var data [100000]byte
var data1 [100000]byte
var data2 [100000]byte
var data3 [100000]byte
var data4 [100000]byte
var data5 [100000]byte
var data6 [100000]byte
var data7 [100000]byte
var data8 [100000]byte
var data9 [100000]byte
var data10 [100000]byte

var benchmarks []BM_DATASIZE

func init() {
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})


	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: [10]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 100, data: [100]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 1000, data: [1000]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10000, data: [10000]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 100000, data: [100000]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 1000000, data: [1000000]byte{}})
}

func BenchmarkAddSizeTestUsing1(b *testing.B){
	cache, e := lru.New(10000000)
	if e != nil {
		fmt.Printf("cache generate error : %s\n",e)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i:= 0 ; i < b.N ; i++{
		cache.Add(key(i), 1)
	}
}
func BenchmarkAddSizeTestUsing2(b *testing.B){
	cache, e := lru.New(10000000)
	if e != nil {
		fmt.Printf("cache generate error : %s\n",e)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i:= 0 ; i < b.N ; i++{
		cache.Add(key(i), 1)
	}
}
func BenchmarkAddSizeTestUsing3(b *testing.B){
	cache, e := lru.New(10000000)
	if e != nil {
		fmt.Printf("cache generate error : %s\n",e)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i:= 0 ; i < b.N ; i++{
		cache.Add(key(i), 1)
	}
}
func BenchmarkAddSizeTestUsing4(b *testing.B){
	cache, e := lru.New(10000000)
	if e != nil {
		fmt.Printf("cache generate error : %s\n",e)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i:= 0 ; i < b.N ; i++{
		cache.Add(key(i), 1)
	}
}
func BenchmarkAddSizeTestUsingBM(b *testing.B){

	for _, bm := range benchmarks {
		cache, e := lru.New(10000000)
		if e != nil {
			fmt.Printf("cache generate error : %s\n",e)
		}
		//b.ResetTimer()
		testName := fmt.Sprintf("Single Add Test : dataSize %d",bm.datasize)
		b.Run(testName, func(b *testing.B) {
			//b.ReportAllocs()
			for i:= 0 ; i < b.N ; i++{
				cache.Add(key(i), data)
			}
		})
	}
}

func BenchmarkAddSizeTestWithout1(b *testing.B){
		cache, e := lru.New(10000000)
		if e != nil {
			fmt.Printf("cache generate error : %s\n",e)
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i:= 0 ; i < b.N ; i++{
			cache.Add(key(i), data8)
		}
}


func BenchmarkAddSizeTestWithout2(b *testing.B){
	cache, e := lru.New(10000000)
	if e != nil {
		fmt.Printf("cache generate error : %s\n",e)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i:= 0 ; i < b.N ; i++{
		cache.Add(key(i), data7)
	}
}


func BenchmarkAddSizeTestWithout3(b *testing.B){
	cache, e := lru.New(10000000)
	if e != nil {
		fmt.Printf("cache generate error : %s\n",e)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i:= 0 ; i < b.N ; i++{
		cache.Add(key(i), data1)
	}
}
