package perf

import (
	"time"
	"fmt"
	"github.com/ground-x/go-gxplatform/core/types"
	"github.com/ground-x/go-gxplatform/rlp"
)

const testCount = 100000
func TestRLPBench(val interface{}, bytes []byte){
	//var body types.Body
	//rlp.DecodeBytes(bytes, &body)



	testRLPEncode(bytes, val)
	testRLPDecode(val)
}

func testRLPEncode(bytes []byte, val interface{}){

	start := time.Now()
	switch val.(type) {
	case *types.Body :
		s := val.(*types.Body)
		fmt.Printf("RLP Test Start Transaction Len : %d, RLP Size : %d \n", len(s.Transactions), len(bytes))
		for i := 0 ; i < testCount ; i++ {
			rlp.DecodeBytes(bytes, &s)
		}
	case *types.Header :
		s := val.(*types.Header)
		for i := 0 ; i < testCount ; i++ {
			rlp.DecodeBytes(bytes, &s)
		}
	}



	elapsedTime := time.Since(start)
	fmt.Printf("Encode Test %s\n",elapsedTime/testCount)
}

func testRLPDecode(val interface{}){
	start := time.Now()
	for i:= 0 ; i < testCount ; i++ {
		_, e := rlp.EncodeToBytes(val)
		if e != nil {
			fmt.Print("error")
		}
	}
	elapsedTime := time.Since(start)
	fmt.Printf("Decode Test %s\n",elapsedTime/testCount)
}
