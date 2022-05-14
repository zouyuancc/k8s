package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type person struct {
	name string
	age  int
}

func main() {
	v := reflect.ValueOf(person{"steve", 30})
	v1, _ := json.Marshal(v)
	fmt.Println(string(v1))
	//count := v.NumField()
	//for i := 0; i < count; i++ {
	//	g := v.FieldByName()
	//	f := v.Field(i)
	//	switch f.Kind() {
	//	case reflect.String:
	//		fmt.Println(f.String())
	//	case reflect.Int:
	//		fmt.Println(f.Int())
	//	}
	//}
}
