package utils

import (
	"bytes"
	"encoding/json"
)

// DeepCopy 深拷贝
// 要传入两个指针，不要传值
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
