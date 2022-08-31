package kmp

import (
	"fmt"
	"testing"
)

func TestNameFindTest(t *testing.T) {
	m := &NameFinder{
		NamePool: nil,
	}
	res := m.FindName("nihao")
	fmt.Println(res)
}

func BenchmarkKmp(b *testing.B) {

	for i := 0; i < b.N; i++ {

	}
}
