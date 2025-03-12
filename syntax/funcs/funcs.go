package main

import "strings"

func Func() {

}

// Func2 有一个参数
func Func2(a int) {

}

// Func3 有多个参数
func Func3(a int, b string) {

}

// Func4 有多个参数,一种类型
func Func4(a, b string) {

}

// Func5 有返回值
func Func5(a, b int) int {
	return 3 // 有返回值，要保证一定返回
}

// Func6 有多个返回值，要全返回
func Func6(a, b int) (int, int) {
	return 3, 5 // 有返回值，要保证一定返回
}

// Func7 返回值有名字
func Func7() (name string, age int) {
	return "zhang san", 18
}

// Func8 返回值有名字， 要么都有名字， 要么都没名字
func Func8() (name string, age int) {
	name = "zhang san"
	age = 18
	return
}

// Func9 直接返回
func Func9() (name string, age int) {
	// 等价于 "", 0
	// 对应类型的零值
	return
}

func Func10(abc string) (string, int) {
	seg := strings.Split(abc, " ")
	return seg[0], len(seg)
}

func Func11(abc string) (first string, length int) {
	seg := strings.Split(abc, " ")
	first = seg[0]
	length = len(seg)
	return
}
