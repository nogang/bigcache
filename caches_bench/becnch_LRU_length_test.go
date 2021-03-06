package main

import (
	"testing"
	"github.com/hashicorp/golang-lru"
	"fmt"
	"os"
	"log"
	"runtime/pprof"
)

type BM_DATASIZE struct{
	datasize int
	data 	interface{}
}

var data = [100000]byte{99999 : 1}


var benchmarks []BM_DATASIZE

func init() {
	/*
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: make([]byte, 10)})

*/
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10, data: [10]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 100, data: [100]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 1000, data: [1000]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 10000, data: [10000]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 100000, data: [100000]byte{}})
	benchmarks = append(benchmarks, BM_DATASIZE{datasize: 1000000, data: [1000000]byte{}})
}

func BenchmarkAddSizeTestUsingBM(b *testing.B){

	for _, bm := range benchmarks {

		testName := fmt.Sprintf("Single Add Test : dataSize %d",bm.datasize)
		b.Run(testName, func(b *testing.B) {
			cache, e := lru.New(10000000)
			if e != nil {
				fmt.Printf("cache generate error : %s\n",e)
			}
			//b.ReportAllocs()

			f, err := os.Create("./" + testName + "_2nd")
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()
			b.ResetTimer()
			for i:= 0 ; i < b.N ; i++{
				cache.Add(key(i), data)
			}
		})
	}
}
