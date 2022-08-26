package redisclix

import (
	"github.com/golang/protobuf/proto"
)

func Encode(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case proto.Message:
		val, err := proto.Marshal(v)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	return value, nil
}

type ArrayReq []interface{}

func NewArrayReq(size int) ArrayReq {
	return make(ArrayReq, 0, size*2)
}

// Add 若value 的类型是proto message，对其进行序列化
func (m *ArrayReq) Add(k interface{}, v interface{}) error {
	data, err := Encode(v)
	if err != nil {
		return err
	}
	*m = append(*m, k, data)
	return nil
}

func NewMapReqStringString(m map[string]string) ArrayReq {
	h := make([]interface{}, 0, len(m)*2)
	for k, v := range m {
		h = append(h, k, v)
	}
	return h
}

func NewMapReqUint64Bytes(value map[uint64][]byte) ArrayReq {
	var values []interface{}
	for k, v := range value {
		values = append(values, k, v)
	}
	return values
}

type ZAddReq map[float64]interface{}

func NewZAddReq() ZAddReq {
	return make(ZAddReq)
}
func (z ZAddReq) Add(score float64, value interface{}) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	z[score] = data
	return nil
}
func (z ZAddReq) Count() int {
	return len(z)
}
