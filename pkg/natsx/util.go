package natsx

import (
	"strings"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// CombineSubject 组合字符串成subject
func CombineSubject(prefix string, s ...string) string {
	if len(s) == 0 {
		return prefix
	}
	bf := bufPool.Get().(*strings.Builder)
	defer func() {
		bf.Reset()
		bufPool.Put(bf)
	}()
	bf.WriteString(prefix)
	for _, v := range s {
		if v == "" {
			continue
		}
		bf.WriteString(".")
		bf.WriteString(v)
	}
	subject := bf.String()

	return subject
}
