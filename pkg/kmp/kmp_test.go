package kmp

import (
	"fmt"
	"testing"
)

func TestKmp(t *testing.T) {
	s := "Hello, 学院君!"
	p := "学院"
	pos := findSubstring(s, p)
	fmt.Printf("Find \"%s\" at %d in \"%s\"\n", p, pos, s)
}
