package cache

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
)

type KeyGenerator interface {
	// Generate 生成key
	// target – the target instance
	// method – the method being called
	// params – the method parameters
	Generate(target any, method string, params ...any) string
}

type SimpleKeyGenerator struct {
}

func (skg SimpleKeyGenerator) Generate(target any, method string, params ...any) string {
	if target == nil {
		target = "target"
	}

	if method == "" {
		method = "method"
	}

	param := "param"

	if len(params) == 0 {
		return fmt.Sprintf("%v:%v:%v", hash32(target), method, hash32(param))
	}

	return fmt.Sprintf("%v:%v:%d", hash32(target), method, hash32(params))
}

func hash32(a any) uint32 {
	array, err := toByteArray(a)
	if err != nil {
		array = []byte{0}
	}

	hash := fnv.New32a()
	hash.Write(array)
	sum32 := hash.Sum32()

	return sum32
}

func toByteArray(a any) ([]byte, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
