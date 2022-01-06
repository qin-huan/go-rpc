package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Person struct {
	Name string `json:"name"`
}

func main() {
	persons := make([]*Person, 0)
	for i:=0; i<5; i++ {
		persons = append(persons, &Person{
			Name: strconv.Itoa(i),
		})
	}
	fmt.Println(marshal(persons))
	person1 := persons[0]
	person1.Name = "chen"
	fmt.Println(marshal(persons))

	fmt.Println(time.Now().Format(time.RFC3339))
}

func marshal(in interface{}) string {
	bytes, _ := json.Marshal(in)
	return string(bytes)
}
