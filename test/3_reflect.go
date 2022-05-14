package main

import (
	"fmt"
	"reflect"
)

type O struct {
	Name string
}

func GetName(o interface{}) string {
	return GetField(o, "Name").(string)
}

func SetName(o interface{}, value interface{}) {
	SetField(o, "Name", value)
}

func main() {
	//获取字段值
	//第一种
	o1 := new(O)
	o1.Name = "d1"
	fmt.Println(GetName(o1))

	//第二种
	o2 := &O{Name: "d2"} //注意&，不带&不是指针会报错
	fmt.Println(GetName(o2))
}

//GetField 获取字段值
func GetField(object interface{}, field string) interface{} {
	t := reflect.TypeOf(object)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		val := reflect.ValueOf(object).Elem().FieldByName(field)
		if !reflect.DeepEqual(val, reflect.Value{}) {
			return val.Interface()
		}
	}

	return nil
}

//SetField 设置字段值
func SetField(object interface{}, field string, value interface{}) {
	t := reflect.TypeOf(object)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		val := reflect.ValueOf(object).Elem().FieldByName(field)
		if !reflect.DeepEqual(val, reflect.Value{}) {
			val.Set(reflect.ValueOf(value))
		}
	}
}
