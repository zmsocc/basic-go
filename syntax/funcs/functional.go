package main

func Functional4() {
	println("hello Functional4")
}

func Functional5(age int) {
	println(age)
}

var Abc = func() string {
	return "hello"
}

// UseFunctional 可以把方法赋值给变量
func UseFunctional() {
	myFunc := Functional4
	myFunc()
	myFunc5 := Functional5
	myFunc5(18)
	/*Abc = func() string {

	}*/
}

func Functional6() {
	// 新定义了一个方法，赋值给了 fn
	fn := func() string {
		return "hello"
	}

	fn()
}

// Functional8 匿名方法立刻发起调用
func Functional8() {
	// 新定义了一个方法，赋值给了 fn
	fn := func() string {
		return "hello"
	}()
	println(fn)
}

// Functional7 他的意思是我返回一个，返回string无参数的方法
func Functional7() func() string {
	return func() string {
		return "hello"
	}
}
