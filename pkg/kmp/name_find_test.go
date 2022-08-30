package kmp

import (
	"fmt"
	"testing"
)

func TestNameFindTest(t *testing.T) {
	res := FindName([]string{"sunihao", "nihaoyy", "zheshini", "haozai"}, "nihao")
	fmt.Println(res)
}
