package main

import (
	"fmt"
	"reflect"
)

/*
使用反射调用方法
*/

func BrushTeeth(name string) string {
	return name + "在刷牙"
}

func PlayFootball(name string) string {
	return name + "在踢足球"
}

func CallAdd(f func(s string) string) {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return
	}

	argv := make([]reflect.Value, 1)
	argv[0] = reflect.ValueOf("葛诗颖")

	reslut := v.Call(argv)
	fmt.Println(reslut[0].String())
}

func main() {
	CallAdd(BrushTeeth)
	CallAdd(PlayFootball)
}
