package main

import (
	"testing"
	"time"
	"github.com/hashicorp/golang-lru"
	"fmt"
)

func BenchmarkAppendFloat(b *testing.B) {
	benchmarks := []struct {
		name    string
		float   float64
		fmt     byte
		prec    int
		bitSize int
	}{
		{"Decimal", 33909, 'g', -1, 64},
		{"Float", 339.7784, 'g', -1, 64},
		{"Exp", -5.09e75, 'g', -1, 64},
		{"NegExp", -5.11e-95, 'g', -1, 64},
		{"Big", 123456789123456789123456789, 'g', -1, 64},

	}
	//dst := make([]byte, 30)
	for _, bm := range benchmarks {


		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				; //AppendFloat(dst[:0], bm.float, bm.fmt, bm.prec, bm.bitSize);
			}
		})
	}
}

func BenchmarkCacheTest(b *testing.B){
	benchmarks := [] struct{
		name		string
		initSize 	int
	}{
		{"lru",100},
		{"lru",1000000},
	}

	var cache *lru.Cache
	for _, bm := range benchmarks {
		switch bm.name {
		case "lru":
			cache, _= lru.New(bm.initSize)
			for i := 0; i < bm.initSize; i++ {
				cache.Add(key(i), value())
			}
		}
		testName := fmt.Sprintf("%s:%d",bm.name, bm.initSize)
		b.Run(testName, func(b *testing.B) {
			for i:= 0 ; i < b.N ; i++{
				cache.Get(key(i%bm.initSize))
			}

		})

	}

}

func TestTime(t *testing.T) {
	testCases := []struct {
		gmt  string
		loc  string
		want string
	}{
		{"12:31", "Europe/Zuri", "13:31"},     // incorrect location name
		{"12:31", "America/New_York", "7:31"}, // should be 07:31
		{"08:08", "Australia/Sydney", "18:08"},
	}
	for _, tc := range testCases {
		loc, err := time.LoadLocation(tc.loc)
		if err != nil {
			t.Fatalf("could not load location %q", tc.loc)
		}
		gmt, _ := time.Parse("15:04", tc.gmt)
		if got := gmt.In(loc).Format("15:04"); got != tc.want {
			t.Errorf("In(%s, %s) = %s; want %s", tc.gmt, tc.loc, got, tc.want)
		}
	}
}

func benchmarkGet(b *testing.B, initCacheSize int, dataSize int){
	cache, _ := lru.New(initCacheSize )
	for i := 0; i < b.N; i++ {
		cache.Add(key(i), value())
	}
}