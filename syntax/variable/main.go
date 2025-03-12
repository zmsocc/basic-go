package main

import "fmt"

// Global 首字母大写，全局可以访问
var Global = "全局变量"

// internal 首字母小写，只能在包内使用，子包也不能使用
var internal = "包内变量， 私有变量"

func main() {
	var a int = 456 // 局部变量
	println(a)

	var b = 234
	println(b)

	var c uint = 15
	println(c)

	var (
		d int = 456
		e int = 345
	)
	println(d, e)

	f := 3
	println(f)
	g := 'a'
	println(fmt.Sprintf("%c", g))
}
