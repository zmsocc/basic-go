package main

import (
	"fmt"
	"unicode/utf8"
)

func String() {
	// he said: "hello go"
	println("he said: \"hello go\"")
	println(`
可以换行
我再换行
`)
	// 只能直接拼接字符串
	println("hello" + "go")

	// println("hello" + string(123))
	println(fmt.Sprintf("hello d%", 123))

	// 计算长度
	println(len("abc"))
	println(len("你好"))
	println(utf8.RuneCountInString("你好"))
}
