package main

import "fmt"

type A struct {
	Name string
	Age int
	Meta *Meta
}

type Meta struct {
	Name string
	Age int
}

func main() {
	//a := &A {
	//	Name: "a",
	//	Age: 18,
	//	Meta: &Meta{
	//		Name: "a",
	//		Age: 10,
	//	},
	//}
	//b := &A {
	//	Name: "b",
	//	Age: 18,
	//	Meta: &Meta{
	//		Name: "b",
	//		Age: 20,
	//	},
	//}
	//
	//var (
	//	ma = make(map[string]interface{})
	//	mb = make(map[string]interface{})
	//)
	//marshal, err := json.Marshal(a)
	//if err != nil {
	//	panic(err)
	//}
	//if err = json.Unmarshal(marshal, &ma); err != nil {
	//	panic(err)
	//}
	//
	//bytes, err := json.Marshal(b)
	//if err != nil {
	//	panic(err)
	//}
	//if err = json.Unmarshal(bytes, &mb); err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("ma: %v, mb: %v\n", ma, mb)
	//m := merge(ma, mb, 0)
	//fmt.Println("m: ", m)

	//i, err := strconv.ParseInt("null", 10, 64)
	//fmt.Printf("%v, %v", i, err)

	var id = "111111"
	var b = []byte(id)
	b[len(b)-1] = '0'
	fmt.Println(len(b), " ", string(b))
}
