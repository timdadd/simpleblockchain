package node

import (
	"encoding/hex"
	"fmt"
	"simpleblockchain/dao"
	"testing"
)

func TestIsBlockHashValid(t *testing.T) {
	testCases := []struct {
		hexHash string
		want    bool
	}{
		{"000000fa04f816039...a4db586086168edfa", true},
		{"123450fa04f816039...a4db586086168edfa", false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s is %t", tc.hexHash, tc.want), func(t *testing.T) {
			var hash = dao.Hash{}
			// Convert it to raw bytes
			hex.Decode(hash[:], []byte(tc.hexHash))
			// Check the validity of the hash
			if got := dao.IsBlockHashValid(hash); got != tc.want {
				t.Errorf("got %t; want %t", got, tc.want)
			}
		})
	}
}

// If Go, you suffix your files with `_test` and prefix the
// testing functions with "Test".
//
// The first argument it the testing helper `t *testing.T`,
// automatically injected by the Go test compiler.

func BenchmarkIsBlockHashValid(b *testing.B) {
	hexHash := "000000fa04f816039...a4db586086168edfa"
	var hash = dao.Hash{}
	// Convert it to raw bytes
	hex.Decode(hash[:], []byte(hexHash))
	for i := 0; i < b.N; i++ {
		dao.IsBlockHashValid(hash)
	}
}
