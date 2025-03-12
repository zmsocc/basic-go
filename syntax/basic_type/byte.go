package main

import "fmt"

func Byte() {
	var a byte = 'a'
	println(a) // 输出的是a的ASCII码
	println(fmt.Sprintf("%c", a))

	var str string = "this is string"
	var bs []byte = []byte(str)
	println(str, bs)
}
