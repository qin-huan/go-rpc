package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
	"strconv"
)

func merge(src, dst map[string]interface{}, depth int) map[string]interface{} {
	if dst == nil {
		return src
	}
	if depth > 16 {
		return dst
	}
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcIsMap := isMap(srcVal)
			dstMap, dstIsMap := isMap(dstVal)

			if srcIsMap && dstIsMap {
				srcVal = merge(srcMap, dstMap, depth+1)
			}
		}
		dst[key] = srcVal
	}
	return dst
}

func isMap(i interface{}) (map[string]interface{}, bool) {
	val := reflect.ValueOf(i)
	m := make(map[string]interface{})

	if val.Kind() == reflect.Map {
		for _, key := range val.MapKeys() {
			m[key.String()] = val.MapIndex(key).Interface()
		}

		return m, true
	} else if val.Kind() == reflect.Slice {
		for i:=0; i<val.Len(); i++ {
			m[strconv.Itoa(i)] = val.Slice(i,i+1)
		}
	}

	return m, false
}

type Person struct {
	Name string `json:"name"`
}

func main() {
	p := &Person{Name: "James"}
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	var i interface{}
	if err = json.Unmarshal(bytes, &i); err != nil {
		panic(err)
	}
	fmt.Println("data: ", i)
	marshal, err := proto.Marshal(i.(proto.Message))
	fmt.Println(string(marshal))
}